# licenses

A tool for inspecting application dependencies requirements/packages/modules in order to extract their license

## Usage
Currently, only the Pip format is supported 

To examine the packages from a requirements file:

`license examine --lang=py requirements.txt`

Or, to examine the packages directly from pip

`pip freeze | license examine --lang=py -`
