# API Reference 

__Note__: I have have not updated this after login was added. Some endpoints might be wrong/missing.

## Endpoints
| Endpoint               | Method | Description |
|------------------------|--------|-------------|
| `/api/upload`         | `POST` | Uploads a file |
| `/api/download`       | `GET`  | Downloads a file |
| `/api/control`        | `POST` | Executes a control command |
| `/api/status`        | `GET`  | Returns recorder status |
| `/api/add-streamer`   | `POST` | Adds a streamer |
| `/api/get-streamers`  | `GET`  | Retrieves a list of streamers |
| `/api/remove-streamer` | `POST` | Removes a streamer |
---

## **1. File Upload**
### **POST /api/upload**
Uploads a file (only `.txt` files are accepted).

#### **Request**
- **Method**: `POST`
- **Headers**: `Content-Type: multipart/form-data`
- **Form Data**:
  - `file` (Required) - The `.txt` file to be uploaded.

#### **Response**
- **Success (200)**
  ```json
  {
    "message": "File 'example.txt' uploaded successfully"
  }
  ```
- **Errors**
  - `405 Method Not Allowed` if method is not `POST`.
  - `400 Bad Request` if file is missing or invalid.

---

## **2. File Download**
### **GET /api/download**
Downloads a dummy `.txt` file.

#### **Request**
- **Method**: `GET`

#### **Response**
- **Headers**:
  - `Content-Type: application/octet-stream`
  - `Content-Disposition: attachment; filename=export.txt`
- **Body**:
  ```
  This is the content of the exported file.
  ```
- **Errors**
  - `405 Method Not Allowed` if method is not `GET`.

---

## **3. Control API**
### **POST /api/control**
Executes a control command (`start`, `stop`, or `restart`).

#### **Request**
- **Method**: `POST`
- **Headers**: `Content-Type: application/json`
- **Body**:
  ```json
  {
    "command": "start"
  }
  ```

#### **Response**
- **Success (200)**
  ```json
  {
    "message": "Control command 'start' executed"
  }
  ```
- **Errors**
  - `405 Method Not Allowed` if method is not `POST`.
  - `400 Bad Request` if JSON payload is invalid.

---

## **4. API Status**
### **GET /api/status**
Returns a string of the recorder status

#### **Request**
- **Method**: `GET`

#### **Response**
- **Success (200)**
  ```json
  {
    "status": "Running"
  }
  ```
- **Errors**
  - `405 Method Not Allowed` if method is not `GET`.

---

## **5. Add Streamer**
### **POST /api/add-streamer**
Adds a new streamer.

#### **Request**
- **Method**: `POST`
- **Headers**: `Content-Type: application/json`
- **Body**:
  ```json
  {
    "data": "streamer_name"
  }
  ```

#### **Response**
- **Success (200)**
  ```json
  {
    "message": "Streamer added successfully",
    "data": "streamer_name"
  }
  ```
- **Errors**
  - `405 Method Not Allowed` if method is not `POST`.
  - `400 Bad Request` if JSON payload is invalid.

---

## **6. Get Streamers**
### **GET /api/get-streamers**
Retrieves a list of all streamers.

#### **Request**
- **Method**: `GET`

#### **Response**
- **Success (200)**
  ```json
  [
    "streamer1",
    "streamer2"
  ]
  ```
- **Errors**
  - `405 Method Not Allowed` if method is not `GET`.

---

## **7. Remove Streamer**
### **POST /api/remove-streamer**
Removes a streamer from the list.

#### **Request**
- **Method**: `POST`
- **Headers**: `Content-Type: application/json`
- **Body**:
  ```json
  {
    "selected": "streamer_name"
  }
  ```

#### **Response**
- **Success (200)**
  ```json
  {
    "message": "Streamer removed successfully",
    "data": "streamer_name"
  }
  ```
- **Errors**
  - `405 Method Not Allowed` if method is not `POST`.
  - `400 Bad Request` if JSON payload is invalid.

---

