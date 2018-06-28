HOME=$(pwd)
MODULES=$HOME/vendor/github.com/Necroforger/Fantasia/modules

# Generate modules and config
go generate

# Generate module dependencies
cd $MODULES/dashboard
go generate

# Create executeable
cd $HOME
go build