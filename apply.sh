FILES=$( ls -dm ./deployment/*.yaml | tr -d ' \n' )

kubectl apply -f $FILES
