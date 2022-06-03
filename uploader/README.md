# BAMS -> CLAMS Uploader Utility

This utility processes BAMS CSV export files exported using the [BAMS Exporter](https://github.com/mikebharris/BAMS/blob/trunk/ExportAttendees.cbl) and posts the entries as messages
to an SQS queue (_clams-nonprod-attendees-input-queue_ by default).  It makes no differentiation between new and
existing attendees and as such isn't all that efficient.  

## Building it

To build the uploader:

```shell
% go build uploader
```

## Testing it

To test it use the _example-data.csv_ file.  Note you'll need to provide AWS credentials sufficient to write to SQS:

```shell
% ./uploader -csv example-data.csv -sqs clams-nonprod-attendee-input-queue
Reading from example-data.csv and writing to clams-nonprod-attendee-input-queue
Queued message # 1  :  {123456 Cyder Punk anicedrop@riseup.net 29 0 04000000 01234 567 890 Fri 0 0 Milk allergy - but that's not a problem with vegan food :)}
Queued message # 2  :  {612297 Rudy Jenkins rudy@jenkins.co 40 0 00000000 07811671893 Fri 0 1 None}
Queued message # 3  :  {BCDEF1 Ronald Chump r.chump@whitehouse.gov 40 0 00000000 01234 567 890 Fri 0 1 }
Queued message # 4  :  {CDEF12 Josefina Rodriguez josefinar@hackitectura.es 29 0 01000000 01234 567 890 Fri 0 2 }
Queued message # 5  :  {DEF123 Random Guy somebody@somewhere.net 29 0 00000000 01234 567 890 Wed 0 0 }
Queued message # 6  :  {EDCBAF Zak Mindwarp zak@mindwarp.io 40 40 01000000 01234 567 890 Wed 0 0 I eat anything}
Queued message # 7  :  {EF1234 Undercover Agent obviouscrusty@gmail.com 29 0 03000000 01234 567 890 Fri 0 0 }
Queued message # 8  :  {FEDCBA Zak Mindwarp zak@mindwarp.io 40 40 01000000 01234 567 890 Wed 0 0 I eat anything}
```

## Using it with BAMS

To use it with BAMS copy the binary file to the same directory as your BAMS executable:

```shell
% cp uploader ~/path/to/BAMS/
```

Now go to BAMS, if you need to build all the binary files with:

```shell
% cd ~/path/to/BAMS/
% ./build.sh
```

Run BAMS (you will need to pass AWS credentials to the environment for the uploader to work within BAMS) and hit the F12 from the Home Screen.
You should see some output like this:

```shell
Reading from attendees.dat and writing to example-data.csv
Total attendees exported to CSV is 001
Reading from  example-data.csv  writing to  clams-nonprod-attendee-input-queue
Queued message # 1  :  {123456 Cyder Punk anicedrop@riseup.net 29 0 04000000 01234 567 890 Fri 0 0 Milk allergy - but that's not a problem with vegan food :)}
.....
```