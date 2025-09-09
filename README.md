## POST /file/{name}
Stores file from `file` post field in files/{name}

## POST /image/{name}
Stores image from `file` post field as webp in files/{name}.webp

## Response
Endpoints respond with `text/plain`.
- Success: status=200 content=OK
- Fail: status!=200 content=ERROR
