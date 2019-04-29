## rs-backup

A simple wrapper around rsync written in Go with mail sending possibilities.

The config is minimalistic, set the flags to true which you would like to be included in the command.

If you set `log` to true, the app will create a temporary file in the system tmp folder and will log the rsync output to that file and will attach this in the email sent at the end.

If the process is complete, there will be an email sent to the specified to adresses in the config. (accepts multiple adresses i.e. array of strings).

If there are any questions or bugs, please feel free to open an issue.