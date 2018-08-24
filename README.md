# BurnQ

BurnQ is very useful demo project ,it's show how to integration with Oracle Cloud API .


* Oracle Cloud Meter API 
   
  Use [Oracle Cloud Meter API](https://docs.oracle.com/en/cloud/get-started/subscriptions-cloud/meter/QuickStart.html) ,it's official public API, you can quick get all about your account  quota information, include usage , unit price ,amount . 
   
* Oracle Jet

* Golang


## How to install 

linux/mac

mkdir brunq

export GOPATH=$PWD

go get -v -u github.com/stevensu1977/burnq

cd src

wget https://github.com/stevensu1977/burnq/raw/master/app.zip

unzip app.zip 


go install github.com/stevensu1977/burnq

## How to Use

./burnq

first step you need init admin account

![Screenshot](https://github.com/stevensu1977/burnq/blob/master/screenshot/init-admin.png?raw=true)

second step you need add already extis Oracle Cloud Account
  ![Screenshot](https://github.com/stevensu1977/burnq/blob/master/screenshot/add-account-01.png?raw=true)

  ![Screenshot](https://github.com/stevensu1977/burnq/blob/master/screenshot/add-account-02.png?raw=true)
    
    
## Screenshot

  ![Screenshot](https://github.com/stevensu1977/burnq/blob/master/screenshot/dashboard.png?raw=true)
  
  
## Licensing
Burn is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/stevensu1977/burnq/blob/master/LICENSE) for the full license text.

Oracle JET is distributed under the [Universal Permissive License(UPL)](https://opensource.org/licenses/UPL).