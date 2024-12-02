package constants

const ECHO = "echo"
const EXIT = "exit"
const TYPE = "type"
const CAT = "cat"

const TOKEN_DELIMITER = " "
const LINE_DELIMITER = '\n'

const BUILTIN = "a shell builtin"

var MapCommand2Path = map[string]string{
	ECHO: BUILTIN,
	EXIT: BUILTIN,
	TYPE: BUILTIN,
	CAT:  BUILTIN,
}
