# back-end-golang


## Daftar Isi

- [back-end-golang ](#back-end-golang)
  - [Daftar Isi](#daftar-isi)
  - [Daftar API](#daftar-api)
    - [Registration](#registration)
    - [Users](#users)



## Daftar API

Kumpulan API tentang data dan informasi yang digunakan dalam capstone project


### Registration

|    Method    |    Developer    |     Endpoint         |                                                                         URL                                                                                             |                                                 Status                                                              | Deskripsi                                                                 |  Penggunaan `Authorization`  |
| ------------ | :-------------: | :------------------------: | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :-----------------------------------------------------------------------------------------------------------------: | :----------------------------------------------------------------- | :-------------------------: |
|    POST      |    M. Hanif     |           /login           | [Link](http://ec2-3-26-30-178.ap-southeast-2.compute.amazonaws.com:8088/swagger/index.html#/User/post_login)                               | `200` Ok     `400` Bad Request     `401` Unauthorized      `404` Not found     `500`Internal Server Error       | Endpoint untuk login                                                       |             Tidak           |
|    POST      |     M. Hanif    |           /register        |[Link](http://ec2-3-26-30-178.ap-southeast-2.compute.amazonaws.com:8088/swagger/index.html#/User/post_register)                            | `201` Ok      `400` Bad Request     `401` Unauthorized      `404` Not found     `500`Internal Server Error       | Endpoint untuk registrasi                                                 |             Tidak           |


### Users

|    Method    |    Developer    |          Endpoint         |                                                                         URL                                                                                             |                                                 Status                                                    | Deskripsi                                                                                      |  Penggunaan `Authorization`  |
| ------------ | :-------------: | :------------------------: | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :---------------------------------------------------------------------------------------------------------------------------: | :----------------------------------------------------------------------- | :-------------------------: |
|     GET      |    M. Hanif     |           /user           | [Link](http://ec2-3-26-30-178.ap-southeast-2.compute.amazonaws.com:8088/swagger/index.html#/User/get_user)                                 | `200` Ok    `400` Bad Request   `401` Unauthorized    `404` Not found   `500`Internal Server Error      | Endpoint untuk mendapatkan detail data user                                                    |               Yes           |
|    PATCH      |     M. Hanif    | /user/update-information |[Link](http://ec2-3-26-30-178.ap-southeast-2.compute.amazonaws.com:8088/swagger/index.html#/User/patch_user_update_information)                            | `200` Ok    `400` Bad Request   `401` Unauthorized    `404` Not found   `500`Internal Server Error       | Endpoint untuk mengupdate data gender, birthdate, profile picture          |               Yes           |
  |      PUT      |     M. Hanif    |   /user/update-password  |[Link](http://ec2-3-26-30-178.ap-southeast-2.compute.amazonaws.com:8088/swagger/index.html#/User/put_user_update_password)                                 | `200` Ok    `400` Bad Request   `401` Unauthorized    `404` Not found   `500`Internal Server Error       | Endpoint untuk mengupdate password user                                      |               Yes           |
|      PUT      |     M. Hanif    |   /user/update-profile   |[Link](http://ec2-3-26-30-178.ap-southeast-2.compute.amazonaws.com:8088/swagger/index.html#/User/put_user_update_profile)                                  | `200` Ok    `400` Bad Request   `401` Unauthorized    `404` Not found   `500`Internal Server Error       | Endpoint untuk mengupdate full name, phone number, birthdate, citizen       |               Yes           |
