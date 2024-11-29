#!/bin/bash


curl -X POST https://id.twitch.tv/oauth2/token \
     -d client_id=g385omhszqp3em7kcntmdglh8tf95j \
     -d client_secret=s3bbtscbl2dtoe5yymbqm3x5mu0t39 \
     -d code=k8nb24z349nlya0fdmtzqfwugbg76h \
     -d grant_type=authorization_code \
     -d redirect_uri=http://localhost:8080/v1/auth/callback
