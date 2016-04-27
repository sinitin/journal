# journal
With this journal program you can keep track of who disturbs you and how much.

Requires mysql to be running with user root and password root.

Limitations:
-Will not work if some tables in the database are deleted.

Usage:
journal list		Prints a list of all disturbances you have logged.
journal hitlist		Prints a list of how many minutes each person have disturbed you.
journal total		Prints the total number of minutes you have been disturbed.
journal log		Logs a new disturbance, submit the arguments [name] [duration in minutes] [reason].
			For example: journal log Sven 15 "Wanted food again"
