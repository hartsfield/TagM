# use in n/vim to restart on save:
# :autocmd BufWritePost * silent! !./autoload.sh
#!/bin/bash
HMACSS="eiwojvioejwoivn_testing_oiewnv4f2332f32fedwe2"
pkill multiparty || true
go build -o multiparty >>log.txt 
./multiparty >>log.txt 2>&1 &
echo http://localhost:15111
echo $(date +%s) > .lastsavetime_bolt
