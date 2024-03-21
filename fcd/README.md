NOTE: My wife and I use this utility to encrypt and decrypt files we need to
share between ourselves. We are both windows users and in tech. It is simple
enough to use while allowing us to do it in a secure way. We also use it to
encrypt files that we need to store on a shared drive. The directions below were
written for our use but can be used by anyone. I am also not a seasoned golang,
so if you see a better way to do something, feel free to keep it to yourself ...
jk, feel free to let me know.

Make sure to copy the folder '.mnemonic' to your users folder. Typically
c:\users\<user_name>. You should also be able to see what folder your users
folder is by typing 'echo %userprofile%' in the command prompt (on windows).

```
'echo %userprofile%'
```

If you are using the comman line and are in the same directory as the
'.mnemonic' folder you can use the following command to copy the folder to your
users folder.

```
xcopy .\.mnemonic %userprofile%\.mnemonic /E /I
```

To encrypt a file called 'passwords.txt', a mnemonic key file called
'our_shared_file_key', and a password 'mySecretPassword123!' use the following
command:

```
.\fcd -e -s passwords.txt -t passwords.txt.encrypted -m our_shared_file_key -p mySecretPassword123!
```

To decrypt a file called 'passwords.txt', a mnemonic key file called
'our_shared_file_key', and a password 'mySecretPassword123!' use the following
command:

```
.\fcd -d -s passwords.txt.encrypted -t passwords.txt -m our_shared_file_key -p mySecretPassword123!
```

Don't forget to delete the original file after encrypting the original file.
Additionally, make sure the original file is on a network or shared drive when
running the encryption commands. This will ensure that the original file is not
visible before the encryption process, and not recoverable after the encryption
process.
