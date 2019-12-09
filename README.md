# Cryptopad Server

## Overview

Cryptopad is a trustless security notepad application for web, desktop, and mobile. It allows users to securely and anonymously have pads stored online without any fear of the contents being compromised and read by anyone else.

## Functionality

A user can create a pad on any client by providing a name (which has never been used on a previous pad), they can then store any data inside said pad such as plain text, formatted text, lists, etc. Once they have added all of the contents they like inside the pad, the contents will be encrypted with a cryptographic key on the client, the encrypted contents will be sent to the server where they will be saved. At any time a client can then download and read a pad, then update the contents provided they have the cryptographic key.

As the server never receives the unencrypted document or cryptographic key, the server could not read the contents even if under legal pressure or if the database was compromised. This is known as trustless security and is the ultimate goal of this application.

## Technical Overview

For creation, a user will have to think of a unique ID for the pad. The client will send a request to the server to check if this ID is available, if it is taken, the user must start again. If it is not taken, the client will generate a random string known as proof, encrypt the empty pad alongside this proof, then send this encrypted pad and proof to the server alongside the proof in plain text.

For loading an existing pad, a user will provide the ID of the existing pad. The client will send a request to the server to download the pad. If it exists, the client will now have to decrypt the pad with the cryptographic key. If all of this is successful, the user can now read the pad and has the proof.

For updating, a client should already have the pad decrypted in this stage, which includes the proof. Now the client can just send the new encrypted contents of this pad, alongside the proof in plain text. The server will now check if the proofs match, if they do, the server will update the contents of the pad in the database to the new contents.

For updating proof, if a user would like to use a new cryptographic key for any reason, the client can send a new proof while updating. If the update is successful, the server will update the proof in the database to the new proof.
For deletion, a user must provide the proof and ID of the pad, which should already be obtained by this point.


## Cryptography Overview

[PBKDF2](https://en.wikipedia.org/wiki/PBKDF2) will be used for generating a 256 bit key from a password.

[AES-256](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) will be used to encrypt the contents with the [CBC cipher mode](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation) using the key generated by [PBKDF2](https://en.wikipedia.org/wiki/PBKDF2).
