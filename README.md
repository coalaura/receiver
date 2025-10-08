## `POST /file/{name}`
Stores file from `file` post field in `files/{name}`

## `POST /image/{name}`
Stores image from `file` post field as webp in `files/{name}.webp`

## `POST /image/{name}/{size}`
Stores image from `file` post field as webp in `files/{name}.webp` with a maximum height/width of `{size}` (128 - 4096)

## Response
Endpoints respond with `text/plain`.
- Success: status=200, content=`OK`
- Fail: status!=200, content=`ERROR`
