# use in n/vim to restart on save:
# :autocmd BufWritePost * silent! !./autoload.sh
#!/bin/bash
HMACSS="eiwojvioejwoivn_testing_oiewnv4f2332f32fedwe2"
pkill tagmachine.xyz || true
go build -o tagmachine.xyz >>log.txt 
./tagmachine.xyz >>log.txt 2>&1 &
echo http://localhost:15111
echo $(date +%s) > .lastsavetime_bolt
