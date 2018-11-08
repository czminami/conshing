# conshing
Conshing is one implementation of Consistent Hashing, and base on [MurmurHash3](https://github.com/spaolacci/murmur3) to apply better Hash Balance.

A major feature of conshing is that, when lookup from conshing, some nodes can excepted and will not return, this feature suitable for service downgrade or server device maintenance.

How to use, please check conshing_test.