# licenses

A tool for inspecting Python requirements in order to extract their license

## Usage
Currently, only the Pip format is supported 

To examine the packages from a requirements file:

`license examine requirements.txt`

Or, to examine the packages directly from pip

`pip freeze | license examine -`
