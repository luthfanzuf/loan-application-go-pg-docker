# loan-application-go-pg-docker

#### Menggunakan Docker-Compose

`docker-compose up -d --build`

#### Fitur Yang Tersedia:
`/user/account/create` customer bisa membuat akun/profil baru <br>
`/user/account/update` customer bisa memperbarui akun/profil <br>
`/user/account/delete` customer bisa menghapus akun/profil <br>
`/user/loan/create` customer bisa membuat pengajuan pinjaman <br>
`/user/loan/update` customer bisa mengedit pengajuan pinjaman yang pernah dibuat dan belum disetujui <br>
`/user/loan/delete` customer bisa menghapus pengajuan pinjaman yang pernah dibuat dan belum disetujui <br>
`/user/account/listall` employee bisa melihat semua list akun yang terdaftar <br>
`/user/loan/listall` employee bisa melihat semua pengajuan yang pernah dibuat <br>

#### Testing bisa mengguankan file `request.rest` yang memanfaatkan [ekstensi REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) di VS Code. <br>


