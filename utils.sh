#!/bin/bash

# insert --- in first line
# for f in `cat failed`
# do
#     echo "--> $f"
#     sed -i '' '1i\
#     ---
#     ' $f
# done

for f in `ls post`
do
    sed -i '' 's/(\/images/(@media/g' post/$f
done
