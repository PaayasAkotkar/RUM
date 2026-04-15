@echo off
echo [LOG] Starting Docker infrastructure (Milvus, etcd, MinIO)...
docker compose up -d

echo [LOG] Building RUM service...
go build -o app.exe .

if %ERRORLEVEL% EQU 0 (
    echo [LOG] Build successful. Launching RUM...
    echo.
    .\app.exe
) else (
    echo [ERROR] Build failed. Please check your code.
    pause
)
