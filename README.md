# lif-notifier

## Example Usage:
`tractrak.exe`

- Download the .exe file to the folder you will have `.lif` files created
- Download the tractrak.key file from your Meet page on the website
- Open a Windows Command Prompt
- `cd` to the folder location
- Run the above command, this message means it's all set
  - ```
     2018/01/18 20:12:52 version: 2.2.0
     2018/01/18 20:12:52 Got the Key
     2018/01/18 20:12:52 Ready to rock this meet ...```
- Successful files copied will output a message: 
  - ```
     2018/01/18 20:29:08 created/modified file: 023-1-01.lif
     2018/01/18 20:29:08  ... trying to upload: 023-1-01.lif
     2018/01/18 20:29:09  ... upload success```
- Hit CTRL-C to end the program.

You may see a message like this:

`2018/01/18 20:29:08  ... read file failed (try again in half-a-second):  open 023-1-01.lif: The process cannot access the file because it is being used by another process.`

Unless it never completes, it will keep trying to upload.

If you see a message like this:

`2018/01/18 20:26:15  ... upload failed:  bad status: 500 Internal Server Error`
 
 (The number and following message may change) the file failed, and will NOT be reattempted. You may re-save or other remove and replace the .LIF file and it will then try again.

## Build ##
- `goreleaser release --snapshot --rm-dist`
- `mv ./dist/lif-copier_windows_386/lif-copier.exe ./tractrak.exe`
- `git tag -a v0.1.0 -m "First release"`
- `git push origin v0.1.0`
