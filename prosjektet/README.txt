In order to run code chmod +x the hall_request_assigner if you are using linux, start elevatorserver and run main.go
Make sure port in elevOperation.go is 15657. This is located in the Initialize function, line 265

Clear txts in DRStorage if you want but is not necessary.
exe files are removed so it will not work on windows without dowloading hall_request_assigner.exe and placing it in PRAssigner folder