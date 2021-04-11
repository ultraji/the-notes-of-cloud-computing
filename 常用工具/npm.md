## npm install -g ... 报错 Error: EACCES: permission denied, access '/usr/local/lib'

Make a directory for global installations:

1. mkdir ~/.npm-global

Configure npm to use the new directory path:

2.npm config set prefix '~/.npm-global'

Open or create a ~/.profile file and add this line:

3.export PATH=~/.npm-global/bin:$PATH

Back on the command line, update your system variables:

4.source ~/.profile