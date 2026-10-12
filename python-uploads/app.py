import os
from flask import Flask, request, jsonify
from werkzeug.utils import secure_filename

app = Flask(__name__)
UPLOAD_DIR = "/data"
os.makedirs(UPLOAD_DIR, exist_ok=True)

@app.get("/")
def home():
    return "XYZ Corporation - Persistent File Upload Service"

@app.post("/upload")
def upload():
    if "file" not in request.files:
        return jsonify(error="No file provided"), 400

    file = request.files["file"]
    filename = secure_filename(file.filename or "")

    if not filename:
        return jsonify(error="Invalid filename"), 400

    file.save(os.path.join(UPLOAD_DIR, filename))
    return jsonify(message="Upload successful", filename=filename), 201

@app.get("/files")
def files():
    return jsonify(files=sorted(os.listdir(UPLOAD_DIR)))

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
