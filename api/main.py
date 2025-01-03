import json
import subprocess
from fastapi import FastAPI, HTTPException, Query, UploadFile, requests
from typing import List, Optional
import os
from fastapi.responses import FileResponse
from pydantic import BaseModel
import requests

app = FastAPI()
final_paths = []


# Model to handle authentication requests
class AuthRequest(BaseModel):
    role_id: str = Query(..., description="Role ID used for Vault authentication")
    secret_id: str = Query(..., description="Role Secret used for Vault authentication")


# Generate a Vault token using AppRole authentication
@app.post("/generate-token")
async def generate_token(
    auth_request: AuthRequest,
    vault_url: str = Query(..., description="Vault URL for authentication"),
):
    """
    Generates a Vault token by authenticating with the AppRole credentials provided.
    """
    data = {
        "role_id": auth_request.role_id,
        "secret_id": auth_request.secret_id,
    }

    response = requests.post(f"{vault_url}/v1/auth/approle/login", json=data)
    if response.status_code == 200:
        token = response.json()["auth"]["client_token"]
        return {"token": token}
    else:
        raise HTTPException(status_code=response.status_code, detail="Authentication with Vault failed")
    
def execute_medusa_command1():
    """
    Executes a Medusa command using environment variables to export secrets from Vault 
    and returns the parsed JSON output.
    """
    try:
        command = [
            "./medusa", "export", os.environ["PATH"],
            "--address", os.environ["ADDRESS"],
            "--token", os.environ["TOKEN"],
            "--format", "json"
        ]
        result = subprocess.run(command, capture_output=True, text=True, check=True)
        return json.loads(result.stdout)
    except subprocess.CalledProcessError as e:
        raise HTTPException(status_code=500, detail=f"Command failed: {e.stderr}")
    except json.JSONDecodeError:
        raise HTTPException(status_code=500, detail="Failed to parse JSON output from medusa")

def execute_medusa_command(address, token, path):
    """
    Executes a Medusa command to export secrets from a specific Vault path and 
    returns the parsed JSON output.
    """
    try:
        command = [
            "./medusa", "export", path,
            "--address", address,
            "--token", token,
            "--format", "json"
        ]
        result = subprocess.run(command, capture_output=True, text=True, check=True)
        return json.loads(result.stdout)
    except subprocess.CalledProcessError as e:
        raise HTTPException(status_code=500, detail=f"Command failed: {e.stderr}")
    except json.JSONDecodeError:
        raise HTTPException(status_code=500, detail="Failed to parse JSON output from medusa")

def generate_paths(data, prefix=""):
    """
    Recursively generates paths from a nested dictionary structure of secrets.
    """
    paths = []
    for key, value in data.items():
        if not isinstance(value, dict):
            continue
        current_path = f"{prefix}/{key}" if prefix else key
        paths.append(current_path)
        paths.extend(generate_paths(value, current_path))
    return paths

@app.get("/get_paths")
async def get_paths(    
    address: str = Query(..., description="Vault URL"),
    token: str = Query(..., description="Vault Token"),
    secret_root: str = Query(..., description="Root secret")
):
    """
    This endpoint retrieves secret paths from Vault based on the provided Vault URL, token, and root secret.

    - **address**: The URL of the Vault server where secrets are stored.
    - **token**: The authentication token required to access Vault.
    - **secret_root**: The root path for fetching secrets.
    """
    global final_paths 
    try:
        secret_paths = execute_medusa_command(address, token, secret_root)  
        final_paths = generate_paths(secret_paths)

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to fetch secret paths: {str(e)}")

    return {"message": "Secret paths fetched successfully", "available path secrets": final_paths}

@app.post("/import")
async def import_to_vault(
    address: str = Query(..., description="Vault URL"),
    token: str = Query(..., description="Vault Token"),
    path: str = Query(..., description="Secret path to import"),
    engine_type: str = Query("kv2", description="Specify the Vault engine type", enum=("kv2", "kv1")),
    files: List[UploadFile] = []
):
    """
    Imports files to a specified Vault path using Medusa. The method accepts multiple files, 
    and the files are temporarily saved to the local system before being uploaded to Vault.

    Returns:
    - A message confirming successful import along with the list of imported file paths
    """
    if not files:
        raise HTTPException(status_code=400, detail="No files provided for import.")

    local_file_paths = []
    try:
        # Save uploaded files temporarily
        for file in files:
            file_location = f"/tmp/{file.filename}"
            with open(file_location, "wb") as f:
                content = await file.read()
                f.write(content)
            local_file_paths.append(file_location)

        # Construct the Medusa command to import the files to Vault
        command = [
            "./medusa", "import", path,
            "--address", address,
            "--token", token,
            "--engine-type", engine_type,
            "--insecure"
        ]
        command.extend(local_file_paths)

        # Run the command and check for errors
        result = subprocess.run(command, capture_output=True, text=True)
        if result.returncode != 0:
            raise HTTPException(status_code=500, detail=f"Medusa import failed: {result.stderr}")

        return {"message": "Files imported successfully", "files": local_file_paths}
    finally:
        # Clean up temporary files
        for file_path in local_file_paths:
            if os.path.exists(file_path):
                os.remove(file_path)

@app.post("/export")
async def export_from_vault(
    address: str = Query(..., description="Vault URL"),
    token: str = Query(..., description="Vault Token"),
    path: str = Query(..., description="Secret path to export"),
    file_name: Optional[str] = Query(None, description="Choose file name for the exported file"),
    output_format: str = Query("yaml", description="Choose the output file format", enum=("yaml", "json")),
    engine_type: str = Query("kv2", description="Specify Vault engine type", enum=("kv2", "kv1"))
):
    """
    Exports secrets from Vault to a file using Medusa. The method allows for optionally,
    choosing the file format (YAML or JSON), and specifying the Vault engine type (kv1 or kv2). The exported file
    will either use the provided filename or a default one generated from the secret path. 

    Parameters:
    - address: Vault URL
    - token: Vault Token for authentication
    - path: Secret path to export from Vault
    - file_name: Optional custom file name for the exported file (defaults to the path name)
    - output_format: The format of the exported file (YAML or JSON, default is YAML)
    - engine_type: Vault engine type (kv1 or kv2, default is kv2)

    Returns:
    - A file response with the exported secret file in the specified format.
    """
    # If file_name is not provided, create a default name based on the path
    file_name_to_use = file_name if file_name else path.replace("/", "_")
    file_name_to_use = f"{file_name_to_use}.{output_format}"

    command = [
        "./medusa", "export", path,
        "--address", address,
        "--token", token,
        "--output", file_name_to_use,
        "--format", output_format,
        "--engine-type", engine_type,
        "--insecure"
    ]
    command = [arg for arg in command if arg]

    # Run the Medusa export command
    result = subprocess.run(command, capture_output=True, text=True)
    if result.returncode != 0:
        raise HTTPException(status_code=500, detail=f"Medusa export failed: {result.stderr}")

    # Check if the file was created
    if not os.path.exists(file_name_to_use):
        raise HTTPException(status_code=500, detail="File was not created.")

    # Return the exported file as a response
    return FileResponse(
        path=file_name_to_use,
        filename=file_name_to_use,
        media_type="application/octet-stream"
    )

@app.post("/copy")
async def copy_secret(
    address: str = Query(..., description="Vault URL"),
    token: str = Query(..., description="Vault Token"),
    source_path: str = Query(..., description="Source secret path in Vault"),
    target_path: str = Query(..., description="Target secret path in Vault"),
    engine_type: str = Query("kv2", description="Specify Vault engine type", enum=("kv2", "kv1"))
):
    """
    This endpoint allows copying secrets from one Vault path (source_path) to another (target_path). 
    The method communicates with Vault using the provided Vault URL and token. It uses the Medusa tool to perform 
    the copy operation. The engine type (kv1 or kv2) for Vault can also be specified. Upon successful execution, 
    a confirmation message with source and target paths is returned. 

    Parameters:
    - address: The URL of the Vault instance.
    - token: The Vault token for authentication.
    - source_path: The path to the secret in Vault to be copied.
    - target_path: The destination path in Vault where the secret will be copied.
    - engine_type: The Vault engine type (kv1 or kv2).

    Returns:
    - A success message indicating the paths of the copied secret.

    Raises:
    - HTTPException: If the copy operation fails or encounters an error.
    """
    try:
        command_copy = [
            "./medusa", "copy", source_path, target_path,
            "--address", address,
            "--token", token,
            "--engine-type", engine_type
        ]
        copy_result = subprocess.run(command_copy, capture_output=True, text=True)
        print(copy_result)

        # Check if the copy operation succeeded
        if copy_result.returncode != 0:
            raise HTTPException(
                status_code=500,
                detail=f"Failed to copy secret: {copy_result.stderr}"
            )
        
        # Return success response with source and target paths
        return {
            "message": f"Secrets from '{source_path}' copied successfully",
            "source_path": source_path,
            "target_path": target_path,
        }
    
    except Exception as e:
        # Handle any other errors during execution
        raise HTTPException(status_code=500, detail=f"Failed to copy secret: {str(e)}")

@app.post("/move")
async def move_secret(
    address: str = Query(..., description="Vault URL"),
    token: str = Query(..., description="Vault Token"),
    source_path: str = Query(..., description="Source secret path in Vault"),
    target_path: str = Query(..., description="Target secret path in Vault"),
    engine_type: str = Query("kv2", description="Specify Vault engine type", enum=("kv2", "kv1"))
):
    """
    This endpoint allows moving secrets from one Vault path (source_path) to another (target_path). 
    The method communicates with Vault using the provided Vault URL and token. It uses the Medusa tool to perform 
    the move operation. After the secret is successfully moved to the target path, it is deleted from the source path. 
    The engine type (kv1 or kv2) for Vault can also be specified. Upon successful execution, 
    a confirmation message with source and target paths is returned. 

    Parameters:
    - address: The URL of the Vault instance.
    - token: The Vault token for authentication.
    - source_path: The path to the secret in Vault to be moved.
    - target_path: The destination path in Vault where the secret will be moved.
    - engine_type: The Vault engine type (kv1 or kv2).

    Returns:
    - A success message indicating the paths of the moved secret.

    Raises:
    - HTTPException: If the move operation fails or encounters an error.
    """
    try:
        command_move = [
            "./medusa", "move", source_path, target_path,
            "--address", address,
            "--auto-approve",
            "--token", token,
            "--engine-type", engine_type
        ]
        move_result = subprocess.run(command_move, capture_output=True, text=True)
        print(move_result)

        # Check if the move operation succeeded
        if move_result.returncode != 0:
            raise HTTPException(
                status_code=500,
                detail=f"Failed to move secret: {move_result.stderr}"
            )
        
        # Return success response with source and target paths
        return {
            "message": f"Secrets from '{source_path}' moved successfully",
            "source_path": source_path,
            "target_path": target_path,
        }
    
    except Exception as e:
        # Handle any other errors during execution
        raise HTTPException(status_code=500, detail=f"Failed to move secret: {str(e)}")

@app.post("/from_old_to_new")
async def migrate(
    old_address: str = Query(..., description="Old Vault URL"),
    old_token: str = Query(..., description="Vault Token for the old Vault"),
    old_path: str = Query(..., description="Secret path in the old Vault"),
    old_engine_type: str = Query("kv2", description="Vault engine type for the old Vault", enum=("kv2", "kv1")), 
    new_address: str = Query(..., description="New Vault URL"),
    new_token: str = Query(..., description="Vault Token for the new Vault"),
    new_path: str = Query(..., description="Secret path in the new Vault"),
    engine_type: str = Query("kv2", description="Vault engine type for the new Vault", enum=("kv2", "kv1")),
    decrypt: bool = Query(False, description="Decrypt secrets during migration", enum=(False, True))
):
    """
    This endpoint migrates secrets from one Vault instance (old_address) to another (new_address).
    It uses the Medusa tool to export secrets from the old Vault and then imports them into the new Vault. 
    The migration process ensures that the secrets are securely transferred.

    Parameters:
    - old_address: The URL of the old Vault instance from where secrets will be exported.
    - old_token: The Vault token for authentication with the old Vault.
    - old_path: The path to the secret in the old Vault to be migrated.
    - old_engine_type: The engine type of the old Vault (kv1 or kv2).
    - new_address: The URL of the new Vault instance where secrets will be imported.
    - new_token: The Vault token for authentication with the new Vault.
    - new_path: The path to the secret in the new Vault where it will be stored.
    - engine_type: The engine type of the new Vault (kv1 or kv2).

    Returns:
    - A success message indicating the completion of the migration, including the old and new Vault addresses.

    Raises:
    - HTTPException: If the export or import operation fails, or if any other error occurs during migration.
    """
    try:
        # Command to export secrets from the old Vault
        export_old_command = [
            "./medusa", "export", old_path,
            "--address", old_address,
            "--token", old_token,
            "--engine-type", old_engine_type,
            "--output", "/tmp/export_old.yaml",
            "--insecure"
        ]
        export_result = subprocess.run(export_old_command, capture_output=True, text=True)

        if export_result.returncode != 0:
            raise HTTPException(
                status_code=500,
                detail=f"Failed to export secret from old Vault: {export_result.stderr}"
            )
        
        # File path of the exported secrets
        file_name = "/tmp/export_old.yaml"
        
        # Command to import secrets into the new Vault
        import_new_command = [
            "./medusa", "import", new_path, file_name,
            "--address", new_address,
            "--token", new_token,
            "--engine-type", engine_type,
            "--decrypt" if decrypt else "",
            "--insecure"
        ]
        import_result = subprocess.run(import_new_command, capture_output=True, text=True)
        print(import_result.stdout)

        if import_result.returncode != 0:
            raise HTTPException(
                status_code=500,
                detail=f"Failed to import secret into the new Vault: {import_result.stderr}"
            )
        
        # Return a success response with old and new Vault addresses
        return {
            "message": f"Secrets from '{old_address}' migrated successfully to '{new_address}'."
        }
    
    except Exception as e:
        # Handle any unexpected errors
        raise HTTPException(status_code=500, detail=f"Failed to migrate secrets: {str(e)}")


@app.delete("/delete")
async def delete_secret(
    address: str = Query(..., description="Vault URL"),
    token: str = Query(..., description="Vault Token"),
    secret_path: str = Query(..., description="Secret path")
):
    """
    This endpoint deletes a secret from a HashiCorp Vault instance.

    It uses the Medusa CLI tool to perform the deletion of the specified secret at the given path in the Vault. 
    The operation is automatically approved and executed with the provided Vault address and authentication token.

    Parameters:
    - address: The URL of the Vault instance where the secret is stored.
    - token: The Vault token used for authentication with the Vault instance.
    - secret_path: The path of the secret to be deleted from the Vault.

    Returns:
    - A success message indicating that the secret at the specified path has been successfully deleted.

    Raises:
    - HTTPException: If the deletion operation fails or if any other error occurs during execution.
    """
    try:
        # Command to delete the secret using Medusa CLI
        command_delete = [
            "./medusa", "delete", secret_path,
            "--address", address,
            "--token", token,
            "--auto-approve",
            "--insecure"
        ]

        # Execute the deletion command
        delete_result = subprocess.run(
            command_delete,
            capture_output=True,
            text=True
        )

        # Check if the command was successful
        if delete_result.returncode != 0:
            raise HTTPException(
                status_code=500,
                detail=f"Failed to delete secret: {delete_result.stderr}"
            )

        # Return a success response
        return {"message": "Secret deleted successfully", "secret_path": secret_path}
    
    except Exception as e:
        # Handle any unexpected errors
        raise HTTPException(status_code=500, detail=f"Failed to delete secret: {str(e)}")


@app.post("/env")
async def export_env(
    address: str = Query(..., description="Vault URL"),
    token: str = Query(..., description="Vault Token"),
    path: str = Query(description="Secret path"),
    engine_type: str = Query("kv2", description="Choose file format", enum=("kv2", "kv1"))
):
    """
    This endpoint sets environment variables for interacting with a HashiCorp Vault instance.

    Parameters:
    - address: The Vault URL to connect to.
    - token: The Vault token used for authentication.
    - path: The base path of the secrets in Vault.
    - engine_type: The type of engine (kv1 or kv2) used in the Vault.

    The endpoint uses these parameters to generate environment variables (`ADDRESS`, `TOKEN`, 
    `ENGINE_TYPE`, and `PATH`) for subsequent Vault operations. It then executes helper functions 
    (`execute_medusa_command1` and `generate_paths`) to retrieve and generate secret paths.

    Returns:
    - A success message indicating the environment variables have been set.
    - A list of available secret paths.
    """
    global final_paths
    try:
        os.environ["ADDRESS"] = address
        os.environ["TOKEN"] = token
        os.environ["ENGINE_TYPE"] = engine_type
        os.environ["PATH"] = path

        secret_paths = execute_medusa_command1()
        final_paths = generate_paths(secret_paths)

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to process ENV variables: {str(e)}")

    return {"message": "ENV variables created successfully", "available path secrets": final_paths}


@app.post("/advenced_export")
async def export_from_vault(
    file_name: Optional[str] = Query(None, description="Choose file name"),
    output_format: str = Query("yaml", description="Choose file format", enum=("yaml", "json")),
    secret_path: str = Query(description="Secret path"),
    engine_type: str = Query("kv2", description="Choose file format", enum=("kv2", "kv1")),
):
    """
    This endpoint exports secrets from a HashiCorp Vault instance to a file.

    Parameters:
    - file_name: The name of the file to save the secrets. Defaults to a name derived from the secret path.
    - output_format: The format of the output file (yaml or json).
    - secret_path: The path of the secret in Vault to export.
    - engine_type: The type of engine (kv1 or kv2) used in Vault.

    Returns:
    - The exported file as a downloadable response.

    Raises:
    - HTTPException: If the Medusa CLI export command fails or the file is not created.
    """
    file_name_to_use = file_name if file_name else secret_path.replace("/", "_")
    file_name_to_use = f"{file_name_to_use}.{output_format}"

    command = [
        "./medusa", "export", secret_path,
        "--address", os.environ["ADDRESS"],
        "--token", os.environ["TOKEN"],
        "--output", file_name_to_use,
        "--format", output_format,
        "--engine-type", engine_type,
        "--insecure"
    ]
    command = [arg for arg in command if arg]

    result = subprocess.run(command, capture_output=True, text=True)
    if result.returncode != 0:
        raise HTTPException(status_code=500, detail=f"Medusa export failed: {result.stderr}")

    if not os.path.exists(file_name_to_use):
        raise HTTPException(status_code=500, detail="File was not created.")
    
    return FileResponse(
        path=file_name_to_use,
        filename=file_name_to_use,
        media_type="application/octet-stream"
    )


@app.post("/advanced_import")
async def import_from_vault(
    path: str = Query(..., description="Secret path to import"),
    files: List[UploadFile] = []
):
    """
    This endpoint imports secrets into a HashiCorp Vault instance from one or more files.

    Parameters:
    - path: The Vault secret path where the secrets will be imported.
    - files: A list of files containing secrets to import.

    The files are temporarily stored in `/tmp` and then imported into Vault using the Medusa CLI.

    Returns:
    - A success message with the list of imported file paths.

    Raises:
    - HTTPException: If no files are provided, the Medusa CLI import command fails, 
      or other errors occur during the import process.
    """
    if not files:
        raise HTTPException(status_code=400, detail="No files provided for import.")

    local_file_paths = []
    try:
        for file in files:
            file_location = f"/tmp/{file.filename}"
            with open(file_location, "wb") as f:
                content = await file.read()
                f.write(content)
            local_file_paths.append(file_location)

        command = [
            "./medusa", "import", path,
            "--address", os.environ["ADDRESS"],
            "--token", os.environ["TOKEN"],
            "--engine-type", os.environ["ENGINE_TYPE"],
            "--insecure"
        ]
        command.extend(local_file_paths)
        result = subprocess.run(command, capture_output=True, text=True)
        if result.returncode != 0:
            raise HTTPException(status_code=500, detail=f"Medusa import failed: {result.stderr}")

        return {"message": "Files imported successfully", "files": local_file_paths}
    finally:
        for file_path in local_file_paths:
            if os.path.exists(file_path):
                os.remove(file_path)


@app.post("/advanced_copy")
async def copy_from_vault(
    source_path: str = Query(..., description="Source secret path"),
    target_path: str = Query(..., description="Target secret path")
):
    await copy_secret(os.environ["ADDRESS"],os.environ["TOKEN"],source_path,target_path,os.environ["ENGINE_TYPE"])

@app.post("/advanced_move")
async def move_from_vault(
    source_path: str = Query(..., description="Source secret path"),
    target_path: str = Query(..., description="Target secret path")
):
    await move_secret(os.environ["ADDRESS"],os.environ["TOKEN"],source_path,target_path,os.environ["ENGINE_TYPE"])

@app.delete("/advanced_delete")
async def delete_from_vault(
    secret_path: str = Query(..., description="Secret path to delete")
):
    await delete_secret(os.environ["ADDRESS"],os.environ["TOKEN"],secret_path)
