package main

import (
	"math/big"
	"test-preconf/spammer"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

var (
	url     = "http://127.0.0.1:28545"
	chainID = big.NewInt(167001) // mainnet
)

var (
	god = spammer.Account{PrivKey: "bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31"} // 0x8943545177806ED17B9F23F0a21ee5948eCaa776

	alice   = spammer.Account{PrivKey: "39725efee3fb28614de3bacaffe4cc4bd8c436257e2c8bb887c4b5c4be45e76d"} // 0xE25583099BA105D9ec0A67f5Ae86D90e50036425
	bob     = spammer.Account{PrivKey: "53321db7c1e331d93a11a41d16f004d7ff63972ec8ec7c25db329728ceeb1710"} // 0x614561D2d143621E126e87831AEF287678B442b8
	charlie = spammer.Account{PrivKey: "ab63b23eb7941c1251757e24b3d2350d2bc05c3c388d06f8fe6feafefb1e8c70"} // 0xf93Ee4Cf8c6c40b329b0c0626F28333c132CF241
	dave    = spammer.Account{PrivKey: "27515f805127bebad2fb9b183508bdacb8c763da16f54e0678b16e8f28ef3fff"} // 0xAe95d8DA9244C37CaC0a3e16BA966a8e852Bb6D6
	eve     = spammer.Account{PrivKey: "7ff1a4c1d57e5e784d327c4c7651e952350bc271f156afb3d00d20f5ef924856"} // 0x2c57d1CFC6d5f8E4182a56b4cf75421472eBAEa4
	fred    = spammer.Account{PrivKey: "5d2344259f42259f82d2c140aa66102ba89b57b4883ee441a8b312622bd42491"} // 0x802dCbE1B1A97554B4F50DB5119E37E8e7336417
	george  = spammer.Account{PrivKey: "3a91003acaf4c21b3953d94fa4a6db694fa69e5242b2e37be05dd82761058899"} // 0x741bFE4802cE1C4b5b00F9Df2F5f179A1C89171A
)

var (
	a1   = spammer.Account{PrivKey: "9ce89a72629e3ea4e4caad5b41c66036039b3208b0785c01bc1f08b9a465a099"} // "0xeA4645C9Cc58f5E976174a27289F92A06C361D34"
	a2   = spammer.Account{PrivKey: "4e03623e9c2a8424357a1f5088f46be603a42db3e71017bd310451ee2f7758e4"} // "0xfeCFEde76a547fD7FC31C3cf361C7c0c55039231"
	a3   = spammer.Account{PrivKey: "b58d9260cf7011226c9166f495d5dc7954c84caf2a81debfd886362a5ddf04b2"} // "0x3ae889734d0C8912CF5D748c1f62a19d74Be4082"
	a4   = spammer.Account{PrivKey: "231a1e47c15bfab926d91c9c616fb41ebbbf1c48b8dc078e474b7208a37336d3"} // "0x7d81f4Da99BAD6DA0AdC277546F5124c5E0a7c61"
	a5   = spammer.Account{PrivKey: "c90f4f2cf571d570859f193cb68efa55b51dd177b53b65ed2dc8ac412a482495"} // "0x2b5e9177EF16203a196B3F109bE657ef16eEb915"
	a6   = spammer.Account{PrivKey: "0e69582680f0849c31884f0de8b840cf62a733fb8a5c17509ef4a62b8b49177f"} // "0xD57Bf6E67695F6836905651Fdd3a3F4151cf3a3C"
	a7   = spammer.Account{PrivKey: "c8bf42b2781fb64ee832f4e2ef4478a9672092e453b217b4de16eda2b1b324e9"} // "0x500aFeeD2E85F9dFD80Ab8cfA042ce7321b9154A"
	a8   = spammer.Account{PrivKey: "0030c159f6a6eb68a26e51a66916c0afe2726b88255fb3c33ae7fcc113b9e9b0"} // "0x06CDb0eA4287C9045763a65B1453EB3836B6C187"
	a9   = spammer.Account{PrivKey: "afc5d6cf796f30cddb39306a1629ae2e1852c1dd4312247789210b3d98047818"} // "0xF4fB8fb2D876BAf872D84812f40704ba5C2D2418"
	a10  = spammer.Account{PrivKey: "f6cd5d7a83f736ad275b3be788f2e7405d83288d47c26c3ec38d61863efeffdd"} // "0xb3F06c91CB93f9792600770C6a7168bB5f74346C"
	a11  = spammer.Account{PrivKey: "4a9ccfe90c898d25d040eb87f44f152d506f07c8733a35bf58482d9d1933b41e"} // "0xAaA8760c34C743C81CA0cbb236cE7c97D2b1F009"
	a12  = spammer.Account{PrivKey: "580a64677e4231f87569a6f19ac64b3be26662bd718dfabdadb0480271f37fd5"} // "0x9FC549F6B336999a8827EA692108628Aab22D0a9"
	a13  = spammer.Account{PrivKey: "e068dc6bb4b8b007c068867700fc291038acc1295dfd914c815c689b24373418"} // "0xd2a3D6BA4a080B93aB485E0aECd91699E37fCA6C"
	a14  = spammer.Account{PrivKey: "20f63e97e2709e5e4262ab3e207ab54332cffd54eb90d2e16dc16ef7c16d0b60"} // "0x924142A34e2b6FaC1f3C1A5DB755028df5Fe6f82"
	a15  = spammer.Account{PrivKey: "ec8fd4213a30a85cd89da8b454c67ba46535536ff3ed1e32c11daac47bce74db"} // "0x365E2512aE06f5e86Ada7f5fd51267e8d076298B"
	a16  = spammer.Account{PrivKey: "f6be7386506547b235e192de361f490a2fd770e1a22e5943717e184a128f2bce"} // "0x85cF62589B9599530e0698979929506b70f2c799"
	a17  = spammer.Account{PrivKey: "0fc8742b6a363b8baf9252b1d4e0642d8684113b61d55c2eaf591a3ba85c210a"} // "0x6eb85D240757E21ad583EE3AF2b5Db5148e97cfa"
	a18  = spammer.Account{PrivKey: "5a410f6096a0bc252b35e6e23a98dc3ccba180d1c4c1a98061b24b2bcc03b4fb"} // "0xeB8492240Fb6346F818e6f8e3Ce38CA30CB712d9"
	a19  = spammer.Account{PrivKey: "b5561454262df0cf69641a871a703512a100c049190022d85ac2df669a3514aa"} // "0xf13b6a674ee8bA1706712276BF78e916A1EfDe52"
	a20  = spammer.Account{PrivKey: "2e3c8f4f799d971280beabd01c60c520cd82221622debbe4acbb75f0fd1acd39"} // "0xcDCf6e9afe18f6030E09222A543Ff15ae709Df30"
	a21  = spammer.Account{PrivKey: "1325236cd08d46e29274168acb1d023ecc0f1000430e38498fa833b457e1bc95"} // "0xf4aC3d0054F8f122003D21685b6278279D53C37F"
	a22  = spammer.Account{PrivKey: "a5d8ca8c7df9076fe45503c3fa6d8854e180430bab9b2b4931b4f5ceb800a2e8"} // "0x0987d781b8F4cD652Dffb9eF994735D68BFE2433"
	a23  = spammer.Account{PrivKey: "f101c55b38138f7d0f9c18ccf2d99f0a7a845f0694821f02e1124c9c3a58d158"} // "0xbFCc2e954A1Ad6FB1e8e72aB3b649a5941dcC24b"
	a24  = spammer.Account{PrivKey: "7ee64263aee62e61d22977a5b5c949b24e5a7a6b71b3317f449780efb5b11c01"} // "0xCae3975181A8d0591BA57a11F259B11E25114C52"
	a25  = spammer.Account{PrivKey: "6539ccfa97a4f764c463b54798ab52b63c509bce3f54847c9eb502cb4db78355"} // "0xF6cD0ff30136846986968c19640F18dc21B0540C"
	a26  = spammer.Account{PrivKey: "bc436206e2fdd034f38477c14eed3be01a41ff7cdebdab2dce4aed607982f437"} // "0x6A93BFB9Cd80D648Ca297B253a0A4c9dF5E099d9"
	a27  = spammer.Account{PrivKey: "ba893b3f3fb2244cbe09baf8406cd85649150391667ff4ebd7152d0d1904a540"} // "0xf89888fEa0bB145EEf5DcC45687c3F5Be2084a11"
	a28  = spammer.Account{PrivKey: "4c9ec219c0c99cb95a8d2cee4e041fc86f7472d8d78f0a751f1c435536737665"} // "0xbBEd21baDCA69393A5b3D6C7fe93afA74cDA3C1A"
	a29  = spammer.Account{PrivKey: "b62346cfb499c9a299d0ef03103a15053004c4590e058e3999c7e6b0598e012f"} // "0xAE028A08A0D0580fe978D6ef0cAF3ed0e351eA80"
	a30  = spammer.Account{PrivKey: "79cb2fadd0c5d4ce44253f294b735d978b919660cc86e6dd39369679333452c5"} // "0x2A3e48A216ECdd8be3DcBc12c2a4E12cdc38c2c8"
	a31  = spammer.Account{PrivKey: "aa6f023905c10d81c09e8ba2879f516443fc5ee88de935da78659313b9c62e09"} // "0x56109Eb02d5626Ee0eebC559EBF78CBEf5E05430"
	a32  = spammer.Account{PrivKey: "c75feee606e50b49d61aaafc0d9094df5459f4bd20aadf815758dfaa0aaf2d59"} // "0x86aFDF97b2588C581133C2A3949d56A67c12F887"
	a33  = spammer.Account{PrivKey: "e48ea2e3eb23f8467e93ef4041378bea78e2b68cc49dc96ce530c3d4e6de1b2e"} // "0xF1cf844a8Bd7Ff46Ba19889E8d357f536fb7B2DC"
	a34  = spammer.Account{PrivKey: "a0cedecb2ee08155f2f2e3fd01712a1d0f9a3666b7669116909d4fb2ac26b5f5"} // "0x8F6A0113910B81Ba8C76353d9DB8b4fe8ec6b838"
	a35  = spammer.Account{PrivKey: "a5abdb0f8d76624b744e4117abbb096f6b6bee78e71cc22a0297e61b13c95b1f"} // "0x946D00945A9f71d5747eB5B8A92606A56F18b4c2"
	a36  = spammer.Account{PrivKey: "20b351f9c77d8b3a9ad0b69e18501b6fca72feb541a3a6db06c84e6d88ef93ef"} // "0xD16229998532AC1E024Af906c426543502f4b466"
	a37  = spammer.Account{PrivKey: "d1a2bbab546bc84d1c3b99a1b12f4c0699286b5aa50de0f1b430b096d276c015"} // "0x48f470397D4aD6171F63C397720D08B24827bF10"
	a38  = spammer.Account{PrivKey: "29abf054392479d89d51146f9806971329cca79144d8d9e90cef224d89174a1f"} // "0x56cd233311D09567596Ba050159Dd59E3cb164A5"
	a39  = spammer.Account{PrivKey: "b886c22c501809d01eeaf6662d1381d51c11c4e4db75afcf11c72916f7684527"} // "0xf8Faf40e7fb72EC985C3a5c155D06a1a809106df"
	a40  = spammer.Account{PrivKey: "a468cb747bf795a6888740388ce20c2e23acf9fd6f85429b58d62c3799688137"} // "0x7204370ce32433E416a98C21EdBb55cDd66eDFd0"
	a41  = spammer.Account{PrivKey: "3ba0ff608237bd97c88c95ce1e05c11bda0307b7e6f7dc21cac5e6f059e8f576"} // "0x4B5002317293e70666F2CB6983c489996aE74e5d"
	a42  = spammer.Account{PrivKey: "ea3a9c69ed82c130bb9e5c7bbaa0f9be789902998071ec43a043aa695c4a3112"} // "0x18827a72cBDA06D1F8808b9d87034D8F02281576"
	a43  = spammer.Account{PrivKey: "5a683ce647673f4b9153bcbab75de4f616c066977a49350c637f686ef6e82bc0"} // "0x3b4110bBAe8e362b7563C6A3bFA6214e57245D40"
	a44  = spammer.Account{PrivKey: "4934361b41a6f545fa83b823b9f83d7a8bbc40b45cabb9bc75d2df7c13949362"} // "0xDD7eB78A3118d5b404DEAffdDaE553277337C7Ed"
	a45  = spammer.Account{PrivKey: "bb02ffc067853d36a8952461cadf7c807ca55a3eda0727119a2e7432b50d3a28"} // "0x3CB8176063B35e9D458B632fAd065aBd2a5bbDB1"
	a46  = spammer.Account{PrivKey: "c1d5f23399d5e8740e3398ff85d609425d3fd655b63acd35c1623ca8f0de3648"} // "0x6a8381af267b9Fb65Ed393Ad115034F2aA9c9504"
	a47  = spammer.Account{PrivKey: "66a89b9247eb9da2040d7d6eee0020dfa60e8a96880de8cadd30207e42ee0fe4"} // "0xB3B5284065E85CbD691cd94117430344A3CaE8B4"
	a48  = spammer.Account{PrivKey: "2e4790a64505f04a962c2974ae046f7c23210b33902a4fe3c2636ee15336a7ff"} // "0xCA14cadcd88D4e4d970d196cc5cadC7BceDD29eE"
	a49  = spammer.Account{PrivKey: "e0ec013899d7d50430f228e0054a90c850df8e60d9208bcf87538379c1c49ca5"} // "0x0824A2dF6da436C53937d343f881CE2cCfafaf8d"
	a50  = spammer.Account{PrivKey: "9eb518fb779a6d4d2be1c247da3a8e42ff68a6b0b599471cf5f21c458535e757"} // "0x9A06b22668259F3744b449709502302A566a7EF6"
	a51  = spammer.Account{PrivKey: "50f51f79625b682452d922436564765b25ba46dc54d734dbe6fc5b049d0ec6c6"} // "0x856d5983787AD9C124Ff57E061D4C897E9f62eBc"
	a52  = spammer.Account{PrivKey: "53a2a60f8326a0977a72b9caa5242954203589e45537a2c13ee2c8ee76313bdd"} // "0x0c4f57a4E2b61ED0fafFF798061972e8ec04DEa7"
	a53  = spammer.Account{PrivKey: "dba3e397dfc6a99289208280e5ef7fb1f572793a74208b450700c3eb761ed748"} // "0x5bA85668449f5A4c4D14355e21Ff58bf5209d0ff"
	a54  = spammer.Account{PrivKey: "fbdd6a1cb36fecf4911e0e577556cd83d9c365073070ace0db6e07c8247e893d"} // "0x3b66005EEbe07aD7bF874B662A06E024D5EfEb29"
	a55  = spammer.Account{PrivKey: "9ecdfc7244b9f1a339c811537bb6fe9c83301a7c81c2df5f0617ad3b9f9ecd38"} // "0x074025AdCdE55A9218247d1F81e00ae2f22aaCcD"
	a56  = spammer.Account{PrivKey: "574ec42482c46965dcb7f192984e1d3f6c08b46fe2324d700513bfbd01867b49"} // "0x66134b64AA4Be2E8EA302D396b90d9E729aaEaC7"
	a57  = spammer.Account{PrivKey: "8e28454f34a6ae6d515619c19f4abdf993cb18d530008003acc3dafd15f1731f"} // "0x8D61597e76d8dc7798cDA867ae5491E6dfcD0536"
	a58  = spammer.Account{PrivKey: "ff0866fbfed42c47e2603ad675ca774499240f8f04d5da6f5c936bb99b6e2534"} // "0xE7865359cdAc077BbFc7983D829a5F6ed580373f"
	a59  = spammer.Account{PrivKey: "64b0ba0e8a708bba1ad3c74239077ab66d771b4b384504e0ede277cab60bec1e"} // "0x6b109E380E6AFCbd4a2d17562b5a1e7847ffc9e8"
	a60  = spammer.Account{PrivKey: "3e0a79e974b7420ff0559671daf6a701299a2240b79fc26bc99462b871c92786"} // "0x758620491f084efE620fAE36D2dE0f5e8655Ab33"
	a61  = spammer.Account{PrivKey: "cc4b280f9afae94ec5c3edbb5a61029d4e85fdbe3e2e8c4c2ba936a681f94fcb"} // "0xf5e2f998b6Cc93F5E7532d2a558805494f6221cf"
	a62  = spammer.Account{PrivKey: "bb8b57d3d514f3b87dd314adcd8a93c85aba5f64780ca886ca9db52b9d2a6325"} // "0x4979040514f991847ead2dbC5fCEb9EF244c2Bf6"
	a63  = spammer.Account{PrivKey: "a27730381aadae9af69d149539bf25e941434c831b9302151456ba34c3e68d6d"} // "0xbbf183027664d0187A1028EDF22C29B39f128d1d"
	a64  = spammer.Account{PrivKey: "f107697b4b64454acba51c934c0d5bbefd850f6d778250b10187e15c83aed252"} // "0xf3629b45872d08B71F04d48e02e1C34390c2310c"
	a65  = spammer.Account{PrivKey: "c8fcc3c9782c22c98aeb1e76473cb84e29f3f3c6e0666e64fe3dda5362640c0a"} // "0xA962C199E45799D3d32baF78ECC4b0B5850bBC19"
	a66  = spammer.Account{PrivKey: "f7b5157165225a90ff95694676a5b4b63ae6d53c8afd3af917e10f93e9c94f30"} // "0xc46A9a956cfbC2BbF3D51C9b436B87Fa38CC4115"
	a67  = spammer.Account{PrivKey: "407ce2c8947bd628d3a1cfc921e06df2227e96c1a79b308958c67668242d04ad"} // "0x45AAfc72bfC931083e17B8dCF4Df4eA6c3dF1De1"
	a68  = spammer.Account{PrivKey: "4f16c97d5769027057787c7dc2f6152cf9cd693f74c81ee9af6318ed3c3220ea"} // "0x5357d6cE4D3E3E3bf44e2D7845E5e8F5BD684cf0"
	a69  = spammer.Account{PrivKey: "31923690c146dbbe387abe5b978c134ccf939b06858ea08baa80d746ed592092"} // "0x88ebBF3F3900A6Fa71f0bBe2942C6316D06396cC"
	a70  = spammer.Account{PrivKey: "f5fed26e96d234033f8914619316e1fb5ae74867230a55749fab6c18165a3769"} // "0xF68B4B87DfefF60D9a911C9138662760a5ce2512"
	a71  = spammer.Account{PrivKey: "ed3010fd5780c31fd5fbbbf7cba25b3c2896a7c4d29317442ddf2c42984f9118"} // "0x64169DE2e58eDdc952dAD4B517EbEfB5E9104567"
	a72  = spammer.Account{PrivKey: "fdf8c0e2859e7047ede988dee8707217d976501ff8d520b0575a491da86719d7"} // "0x4BC68dE40E38500cFcF5ACB56f042Df6713dd7EE"
	a73  = spammer.Account{PrivKey: "098009ff8cd24d6e69a7af0690a4b5eb70b9786b7396a4a32669cbe3957f8f39"} // "0x564B122e0fD4870dA4E2CA12D0C04DC3F1CC6037"
	a74  = spammer.Account{PrivKey: "e6aa7e3df93f8d229a3c63897f268b5dc028a43c9682ad2f781b013d83d9a056"} // "0x0633F1C655e57311655C907D295BF5Dce9fd4941"
	a75  = spammer.Account{PrivKey: "3d9530f6d644f3fab979a5bc2c8b28912b5d639585679054b1c72add1efb15c1"} // "0xbb7279043B391eE0Eef04A8976F71D9E35199080"
	a76  = spammer.Account{PrivKey: "d86eadf8efabdb8656fd495c32b309d2aaef8735bd073a5ed92128a9c6fdc2e3"} // "0xD791C573783D7fA60F9a5b8F65551178fF11da56"
	a77  = spammer.Account{PrivKey: "7e5cdb29c4eafd9d22f842cb40b43f8f6c2430e311aadc585f2f493034f6fc43"} // "0x3409e543195599D077e6525063f9453cBC948Ab0"
	a78  = spammer.Account{PrivKey: "683d6e8dd04a4fc64cce72db65dafd2e41f0843ba4126c99ec2d6b0b6a719b41"} // "0xBeAbd8f49E509E3AA7975D686B893761f6A3C120"
	a79  = spammer.Account{PrivKey: "4b3ca76223c4f12ca0b54620031f704c34ab9f3a630241b04acf3fcb71b58894"} // "0x17B811C3132B5Eb81730d223Bf789Fdaaf805063"
	a80  = spammer.Account{PrivKey: "e4abef925f5bbac176073e63fef3267b69cb2522748b4f4b492219b4ac5381e8"} // "0x5f6BaFbB479a0d00f87Fa1BCC8D4dBd8e10efb30"
	a81  = spammer.Account{PrivKey: "03bcbd6d5a946aeef61b05d4882415f1bd7ea819ad55e332472e8e9a57c46ce9"} // "0xd79FbA7709B5317dc627c19EE2a6fEf2A283C91f"
	a82  = spammer.Account{PrivKey: "10f45d092b8e9d714319b6d503a1cc12ddef0c3d782c9ce8e2c64bb1cd53c88c"} // "0x6f4a2F1625cBc759DDacafdCaC22Da10875cd60c"
	a83  = spammer.Account{PrivKey: "88235fee7674c2e2a1454967edcadc8fd363a276091e07dbc89f9e64f708697f"} // "0x3523C5eA57aD70aD203110E38E2f30732034Ba01"
	a84  = spammer.Account{PrivKey: "20a50d34f168c9a1b0ac2f84eb5b77150b5d62a1e61b9d2990bd21909552169a"} // "0x7Cb9d9D7312730954dFB86017f2DC4a576F44981"
	a85  = spammer.Account{PrivKey: "29352d0cb826814a4d45116763e78985e99490678b4e3551731a743149bb592e"} // "0xa6cFE46A90ffB5BA74825aa32C0dDCF123869e44"
	a86  = spammer.Account{PrivKey: "1f71de316e6f92914997464bbd19fb674e94e88515f5ca8a293c6f0a406c5fb6"} // "0x15ABf22faa9BB0A0338F66f39761D8dAc8E9bA36"
	a87  = spammer.Account{PrivKey: "c8bc1551146d2c4eca7ea5bf6c22941f4541f21adc09b44e81c681aaff5c48d3"} // "0x4bfE504D3D8FcC1660f770e59f305F60042e6AA4"
	a88  = spammer.Account{PrivKey: "048f0d3289ad4dc7c253afdb83dd76d8aed76f5ccedc5a962624d7fa6d89d4c2"} // "0x97Ff0f1079Fc3e9Da16b73eA7365341141d3E0B2"
	a89  = spammer.Account{PrivKey: "66b3a909a772ce62f3b27ed0c9960bc6039026b631a217dfa56327be8179c679"} // "0x898A316ea345e5Dcbff95111fF7043923c9F8476"
	a90  = spammer.Account{PrivKey: "9a3e52d7c02b6fe3d611a1106294185c1c3a6bf8c13564324fa75eae788e8fef"} // "0xc65B8B0B7A2592C162517134607D9D3587A83410"
	a91  = spammer.Account{PrivKey: "61dbc4920192576130830929ae6079b036ed2109b3f43b6e00c57e41d33362df"} // "0xCe1289598c1820773D00b7f3B745FF973C4606eB"
	a92  = spammer.Account{PrivKey: "99c8fdfb33dbc8e85de40db5182f026fa7c40273d063fba4ed6377276ca1fa8c"} // "0xa2f654b69ee3826Dd4bc62b9c766B54cAeBFDC12"
	a93  = spammer.Account{PrivKey: "3ffeb9fdc23f7571014223bf7b3fdfa0f89db0fe8dd8ff76628981dbda037bee"} // "0xb0D0157aD43A3DAcce43191E8036CCd221ddB63d"
	a94  = spammer.Account{PrivKey: "7ae4d743055fe597743ebc93d67df34adbb6d703980f6a7d0f59829a2dfecba7"} // "0x24213daadFaE155cb7E6f098a1d53f79CEaB5f63"
	a95  = spammer.Account{PrivKey: "f2f069a8bc1fcca83f251ae7190bc95e971f52a578f60cf9c3a53829baf9ddd2"} // "0xA352990093c49cB04eF9fB063B0Ade119279f8E8"
	a96  = spammer.Account{PrivKey: "af82c8ddc456178aec8ef38251484ce3083381383d9d6bf9ed943d6ca1e7ad74"} // "0x20b5De9f8ab7A8e727FC982964841622B4a3B56f"
	a97  = spammer.Account{PrivKey: "b5ce3634df386e5dedc71bc3110c6d8da43e37740c31dcd3444c3a994239b25b"} // "0x41bdca7C8574Ce87AaB9D7a5002c1f06170a8Efa"
	a98  = spammer.Account{PrivKey: "aa9246d2327389e1d7a4551df15ba9cb6858fa87dc4e3dac07ca563bc43c3a87"} // "0x455a1D723080090DC631518446C769F682B19804"
	a99  = spammer.Account{PrivKey: "a7c8ff39289d0109fe07b1cf6c3153e4711d9115d41d64121f757903a81705b3"} // "0x91e5447529F6D31eBE4dd598E0724108C5da5662"
	a100 = spammer.Account{PrivKey: "05a65d0aad96a5770338f75ad3d7b95d4236c296d5b90c45074abc15e43180da"} // "0xBf444a6610ed5A084006E645B8B6116dC8AD8267"
)

var (
	accounts = []*spammer.Account{&alice} // &a1, &a2, &a3, &a4, &a5, &a6, &a7, &a8, &a9, &a10, &a11, &a12, &a13, &a14, &a15, &a16, &a17, &a18, &a19, &a20, &a21, &a22, &a23, &a24, &a25, &a26, &a27, &a28, &a29, &a30, &a31, &a32, &a33, &a34, &a35, &a36, &a37, &a38, &a39, &a40, &a41, &a42, &a43, &a44, &a45, &a46, &a47, &a48, &a49, &a50, &a51, &a52, &a53, &a54, &a55, &a56, &a57, &a58, &a59, &a60, &a61, &a62, &a63, &a64, &a65, &a66, &a67, &a68, &a69, &a70, &a71, &a72, &a73, &a74, &a75, &a76, &a77, &a78, &a79, &a80, &a81, &a82, &a83, &a84, &a85, &a86, &a87, &a88, &a89, &a90, &a91, &a92, &a93, &a94, &a95, &a96, &a97, &a98, &a99, &a100
)

var (
	maxTxsPerAccount = uint64(1)

	gasPriceMultiplier = big.NewInt(1_000)
	defaultValue       = big.NewInt(1_000_000_000_000) // (1 eth = 1_000_000_000_000_000_000 wei)
	defaultGas         = uint64(21_000)
	defaultData        = make([]byte, 0)

	txDefaults = func(i interface{}) *types.Transaction {
		var tx *types.Transaction

		switch txData := i.(type) {
		case types.LegacyTx:
			if txData.Value == nil {
				txData.Value = defaultValue
			}
			txData.Gas = defaultGas
			if txData.GasPrice == nil {
				txData.GasPrice = new(big.Int).Mul(big.NewInt(10), gasPriceMultiplier)
			}
			txData.Data = defaultData
			tx = types.NewTx(&txData)
		case types.AccessListTx:
			if txData.Value == nil {
				txData.Value = defaultValue
			}
			txData.Gas = defaultGas
			if txData.GasPrice == nil {
				txData.GasPrice = new(big.Int).Mul(big.NewInt(7), gasPriceMultiplier)
			}
			txData.Data = defaultData
			tx = types.NewTx(&txData)
		case types.DynamicFeeTx:
			if txData.Value == nil {
				txData.Value = defaultValue
			}
			txData.Gas = defaultGas
			if txData.GasFeeCap == nil {
				txData.GasFeeCap = new(big.Int).Mul(big.NewInt(8), gasPriceMultiplier)
			}
			if txData.GasTipCap == nil {
				txData.GasTipCap = new(big.Int).Mul(big.NewInt(5), gasPriceMultiplier)
			}
			txData.Data = defaultData
			tx = types.NewTx(&txData)
		case types.BlobTx:
			if txData.Value == nil {
				txData.Value = uint256.NewInt(defaultValue.Uint64())
			}
			txData.Gas = defaultGas
			if txData.GasFeeCap == nil {
				txData.GasFeeCap = new(uint256.Int).Mul(uint256.NewInt(9), uint256.NewInt(gasPriceMultiplier.Uint64()))
			}
			if txData.GasTipCap == nil {
				txData.GasTipCap = new(uint256.Int).Mul(uint256.NewInt(7), uint256.NewInt(gasPriceMultiplier.Uint64()))
			}
			txData.Data = defaultData
			tx = types.NewTx(&txData)
		case types.InclusionPreconfirmationTx:
			if txData.Value == nil {
				txData.Value = defaultValue
			}
			txData.Gas = defaultGas
			if txData.GasFeeCap == nil {
				txData.GasFeeCap = new(big.Int).Mul(big.NewInt(1), gasPriceMultiplier)
			}
			if txData.GasTipCap == nil {
				txData.GasTipCap = new(big.Int).Mul(big.NewInt(0), gasPriceMultiplier)
			}
			txData.Data = defaultData
			tx = types.NewTx(&txData)
		}

		return tx
	}
)

var (
	txsPerAccount = func(account *spammer.Account, nonce uint64, currentSlot uint64, assignedSlots []uint64) []interface{} {
		firstAcceptableSlot := spammer.CalculateFirstAcceptableSlot(currentSlot, assignedSlots)

		currentSlotDeadline := new(big.Int).SetUint64(currentSlot)

		pastSlotDeadline := new(big.Int).Add(firstAcceptableSlot, big.NewInt(-1))

		acceptableSlotDeadline1 := new(big.Int).Add(firstAcceptableSlot, big.NewInt(1))
		acceptableSlotDeadline2 := new(big.Int).Add(firstAcceptableSlot, big.NewInt(2))
		acceptableSlotDeadline3 := new(big.Int).Add(firstAcceptableSlot, big.NewInt(3))
		// acceptableSlotDeadline4 := new(big.Int).Add(firstAcceptableSlot, big.NewInt(4))
		// acceptableSlotDeadline5 := new(big.Int).Add(firstAcceptableSlot, big.NewInt(5))
		// acceptableSlotDeadline6 := new(big.Int).Add(firstAcceptableSlot, big.NewInt(6))
		// acceptableSlotDeadline7 := new(big.Int).Add(firstAcceptableSlot, big.NewInt(7))
		// acceptableSlotDeadline8 := new(big.Int).Add(firstAcceptableSlot, big.NewInt(8))
		// acceptableSlotDeadline9 := new(big.Int).Add(firstAcceptableSlot, big.NewInt(9))
		// acceptableSlotDeadline10 := new(big.Int).Add(firstAcceptableSlot, big.NewInt(10))

		// offsetToEpochEnd := big.NewInt(int64(uint64(common.EpochLength) - common.SlotIndex(currentSlotDeadline.Uint64())))
		// endSlotDeadline := new(big.Int).Add(currentSlotDeadline, offsetToEpochEnd)
		// tooFarInFutureSlotDeadline := new(big.Int).Add(endSlotDeadline, big.NewInt(common.SlotsOffsetInAdvance))

		txs := map[common.Address][]interface{}{
			*god.Address(): {
				// types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: alice.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: bob.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: dave.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 4, To: eve.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 5, To: fred.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 6, To: george.Address(), Deadline: firstAcceptableSlot},

				// types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: a1.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: a2.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: a3.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 4, To: a4.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 5, To: a5.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 6, To: a6.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 7, To: a7.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 8, To: a8.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 9, To: a9.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 10, To: a10.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 11, To: a11.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 12, To: a12.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 13, To: a13.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 14, To: a14.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 15, To: a15.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 16, To: a16.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 17, To: a17.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 18, To: a18.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 19, To: a19.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 20, To: a20.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 21, To: a21.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 22, To: a22.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 23, To: a23.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 24, To: a24.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 25, To: a25.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 26, To: a26.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 27, To: a27.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 28, To: a28.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 29, To: a29.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 30, To: a30.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 31, To: a31.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 32, To: a32.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 33, To: a33.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 34, To: a34.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 35, To: a35.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 36, To: a36.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 37, To: a37.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 38, To: a38.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 39, To: a39.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 40, To: a40.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 41, To: a41.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 42, To: a42.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 43, To: a43.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 44, To: a44.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 45, To: a45.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 46, To: a46.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 47, To: a47.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 48, To: a48.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 49, To: a49.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 50, To: a50.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 51, To: a51.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 52, To: a52.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 53, To: a53.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 54, To: a54.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 55, To: a55.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 56, To: a56.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 57, To: a57.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 58, To: a58.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 59, To: a59.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 60, To: a60.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 61, To: a61.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 62, To: a62.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 63, To: a63.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 64, To: a64.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 65, To: a65.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 66, To: a66.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 67, To: a67.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 68, To: a68.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 69, To: a69.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 70, To: a70.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 71, To: a71.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 72, To: a72.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 73, To: a73.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 74, To: a74.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 75, To: a75.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 76, To: a76.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 77, To: a77.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 78, To: a78.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 79, To: a79.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 80, To: a80.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 81, To: a81.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 82, To: a82.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 83, To: a83.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 84, To: a84.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 85, To: a85.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 86, To: a86.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 87, To: a87.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 88, To: a88.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 89, To: a89.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 90, To: a90.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 91, To: a91.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 92, To: a92.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 93, To: a93.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 94, To: a94.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 95, To: a95.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 96, To: a96.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 97, To: a97.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 98, To: a98.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 99, To: a99.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 100, To: a100.Address(), Deadline: firstAcceptableSlot},
			},
			*alice.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 4, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 5, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 6, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 7, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 8, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 9, To: charlie.Address(), Deadline: acceptableSlotDeadline1},

				// types.InclusionPreconfirmationTx{Nonce: nonce + 10, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 11, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 12, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 13, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 14, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 15, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 16, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 17, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 18, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 19, To: charlie.Address(), Deadline: firstAcceptableSlot},

				// types.InclusionPreconfirmationTx{Nonce: nonce + 20, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 21, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 22, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 23, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 24, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 25, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 26, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 27, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 28, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 29, To: charlie.Address(), Deadline: firstAcceptableSlot},

				// types.InclusionPreconfirmationTx{Nonce: nonce + 30, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 31, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 32, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 33, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 34, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 35, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 36, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 37, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 38, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 39, To: charlie.Address(), Deadline: firstAcceptableSlot},

				// types.InclusionPreconfirmationTx{Nonce: nonce + 40, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 41, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 42, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 43, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 44, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 45, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 46, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 47, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 48, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 49, To: charlie.Address(), Deadline: firstAcceptableSlot},

				// types.InclusionPreconfirmationTx{Nonce: nonce + 50, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 51, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 52, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 53, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 54, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 55, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 56, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 57, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 58, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 59, To: charlie.Address(), Deadline: firstAcceptableSlot},

				// types.InclusionPreconfirmationTx{Nonce: nonce + 60, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 61, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 62, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 63, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 64, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 65, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 66, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 67, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 68, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 69, To: charlie.Address(), Deadline: firstAcceptableSlot},

				// types.InclusionPreconfirmationTx{Nonce: nonce + 70, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 71, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 72, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 73, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 74, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 75, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 76, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 77, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 78, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 79, To: charlie.Address(), Deadline: firstAcceptableSlot},

				// types.InclusionPreconfirmationTx{Nonce: nonce + 80, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 81, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 82, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 83, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 84, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 85, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 86, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 87, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 88, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 89, To: charlie.Address(), Deadline: firstAcceptableSlot},

				// types.InclusionPreconfirmationTx{Nonce: nonce + 90, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 91, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 92, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 93, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 94, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 95, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 96, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 97, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 98, To: charlie.Address(), Deadline: firstAcceptableSlot},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 99, To: charlie.Address(), Deadline: firstAcceptableSlot},
			},
			*bob.Address(): {
				types.DynamicFeeTx{Nonce: nonce + 0, To: alice.Address()},
				types.DynamicFeeTx{Nonce: nonce + 1, To: alice.Address()},
				types.DynamicFeeTx{Nonce: nonce + 2, To: alice.Address()},
				types.LegacyTx{Nonce: nonce + 3, To: alice.Address()},
				types.LegacyTx{Nonce: nonce + 4, To: alice.Address()},
				types.LegacyTx{Nonce: nonce + 5, To: alice.Address()},
				types.AccessListTx{Nonce: nonce + 6, To: alice.Address()},
				types.AccessListTx{Nonce: nonce + 7, To: alice.Address()},
				types.AccessListTx{Nonce: nonce + 8, To: alice.Address()},
				types.DynamicFeeTx{Nonce: nonce + 9, To: alice.Address()},
				types.DynamicFeeTx{Nonce: nonce + 10, To: alice.Address()},
				types.DynamicFeeTx{Nonce: nonce + 11, To: alice.Address()},
			},
			*charlie.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 4, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 5, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 6, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 7, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 8, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 9, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.DynamicFeeTx{Nonce: nonce + 6, To: george.Address()}, // does not have immediate receipt
			},
			*dave.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: george.Address(), Deadline: pastSlotDeadline},    // rejected, past slot
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: george.Address(), Deadline: currentSlotDeadline}, // rejected, current slot
			},
			*eve.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 10, To: george.Address(), Deadline: acceptableSlotDeadline1}, // disallowed, preconf with future nonces
				types.InclusionPreconfirmationTx{Nonce: nonce + 11, To: george.Address(), Deadline: acceptableSlotDeadline1}, // disallowed, preconf with future nonces
			},
			*fred.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: george.Address(), Deadline: acceptableSlotDeadline3},
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: george.Address(), Deadline: acceptableSlotDeadline3},
				types.DynamicFeeTx{Nonce: nonce + 2, To: george.Address()}, // processed, due to higher fees
				types.DynamicFeeTx{Nonce: nonce + 3, To: george.Address()}, // processed, due to higher fees
			},
			*george.Address(): {
				types.DynamicFeeTx{Nonce: nonce + 0, To: alice.Address()},
				types.DynamicFeeTx{Nonce: nonce + 1, To: alice.Address()},
				types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: alice.Address(), Deadline: acceptableSlotDeadline1}, // not preconfirmed
				types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: alice.Address(), Deadline: acceptableSlotDeadline1}, // not preconfirmed
			},
		}

		return txs[*account.Address()]
	}
)
