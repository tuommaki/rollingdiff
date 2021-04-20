// Package fastcdc provides an implementation for FastCDC algorithm.
// For original paper, see: https://www.usenix.org/conference/atc16/technical-sessions/presentation/xia
package fastcdc

// Following table is generated with (using crypto/rand):
//
// 	max := big.NewInt(0)
//  max.SetUint64(^uint64(0))
//
//  for i := 1; i <= 256; i++ {
//  	n, err := rand.Int(rand.Reader, max)
//  	if err != nil {
//  		panic(err)
//  	}
//  	fmt.Printf("%#x, ", n)
//
//  	if i % 4 == 0 {
//  		fmt.Printf("\n")
//  	}
//
//  }
//
var gear = [256]uint64{
	0x7e0d7b55dbee256b, 0xad0015907824fa8c, 0xc998baa44ddc38c4, 0x12e843a3af83dc8,
	0x2b548fa27633a361, 0x103ade8af889f00, 0x3f43cdde7f796e78, 0x4d64e80486873605,
	0xbbedba3deb651de4, 0xf37bd7319d68419e, 0x44ec0fcb25d5d7c6, 0xa7cb8488c922bcd3,
	0x1abdbcc3ca4f1f76, 0x9b9d2960c8023d4, 0x9f85ba3711732bc9, 0xdde852b47995194e,
	0x325b91c50d421ce6, 0xed9b011759f7967b, 0x8df87e21a087fe90, 0xb465d6e6581277d8,
	0x47552b145e684832, 0x7fd57cb34b46ceee, 0x34e17a53fd4ba44d, 0xc3423bc54ebbfa31,
	0x348654ed8cbfc9bd, 0xb06ddc1a7955d8f3, 0x201ed103cf444cc1, 0x2c0845e5bf15cf05,
	0x5f73485a220dc7a2, 0x5dd438ceb1b5f5cd, 0x1d0204df28191f95, 0x2533eea6672ea03,
	0x7bb642190c8a9603, 0xdb04930961409c, 0x7d421a95845fb9f4, 0xbf4450ba567257aa,
	0x24bcfc5241b1c1fd, 0xe4c460596f7de2a6, 0x3a975891cc403771, 0xc7a2d68bfae212b9,
	0xf53fc8ad0266c144, 0xf92d1b063dcb126e, 0x3df6cd7b3e6399c4, 0xb3a375548607929e,
	0x929e038a4c5da1c7, 0x58e136b452543a06, 0x6290cb21e5b4adf4, 0x9a4beb2844f0a5de,
	0x7ae9bfb52308e399, 0xdbe4bd89f3649858, 0xbb1500c73d3454e9, 0xa7e687ab1d2ec176,
	0xf81692146d3b244c, 0xd7f21d8160813977, 0x2a6d003fe496b093, 0xeec2304055c22a6a,
	0x8a073b8e958b81a, 0x58d3a95a55d8d342, 0xf507bdb2e261f79c, 0xd49c5d01089e5609,
	0x7c3e4f1b913321bc, 0xac2c3d0d25858936, 0x54b3c89a7fc0e785, 0xe5aedeeda8847d52,
	0xa61501ce157a943b, 0x263d4278daaf3646, 0xd039ea152f56faf6, 0x56eacc93647e4ad9,
	0x4aaf3c558521f08, 0xc07e9abc4c50dfa7, 0x417376596a8b83e4, 0x242477dfcc35666,
	0xd5727958fd71253c, 0x1de7d4b1941dfd04, 0x3a2f7fef33abd118, 0xa52d13a49160ff1,
	0x6c7063eccc1d3553, 0x7227786b7a071bcf, 0x33b7a68f43c8588d, 0xb3859db42f0abacf,
	0xa5f1bd033f351779, 0xd6270e2a76b7a484, 0x8e26f6bd1242fcde, 0xe4c35cd00eb33310,
	0xb39dc5db27b96859, 0xdd042e52c293db93, 0xb9e952c9b753cb41, 0xb42b01ffe1d90379,
	0xc0a4a1c1cebffb88, 0xd30fede7f6ee08a8, 0xfb395586d0f35cd2, 0x641f5f622a30cc65,
	0x1b318488887162a1, 0x9c875b122b724990, 0xb5be7e78fdca1386, 0x995b8a7b289fad95,
	0x695f817733e12d82, 0x5427212f2ebd5433, 0xc5bbf300dfadbdc, 0x327f254227f7845f,
	0xf3302afbc7d29299, 0x17de837fc850a290, 0xec9dd3c7bdb02e80, 0xd6c439ef847c06f1,
	0xcc46fd07b6b5fda8, 0x7f3cf4acdf4bf90f, 0xd3c360c0724c115e, 0xbe2fe645fb40a971,
	0x7e1d01b143589e9d, 0xb4f87faaf08f18a6, 0xe03703f6060acb61, 0x1d32395d42b2e0cd,
	0x12f4625d1d849285, 0x561436e6e2a4ac1f, 0xd0c56eb16e9e99da, 0x58441825156b89b7,
	0x9020e0eea94db0dc, 0x86becccddc649077, 0x3e7dec8aec61b91d, 0x6846fa00eea02e92,
	0xc068b567d20e40f9, 0xef80f835bfec7a09, 0x957a025fb19526bd, 0x6f270baa17ac36e2,
	0xa1025f15c1de32ad, 0xd2d9e05c667964b2, 0xf7990e46fa34d401, 0x66f31732dd947385,
	0xaa893caef85a62c8, 0x92e40306c87d1317, 0xb69839eb92ac0b97, 0xa8385273ab1005bf,
	0xd27c90f9f0188645, 0x24ae10cb80f2b2b9, 0xcb572088e18be6bf, 0x996e7c43c30855c4,
	0xedaa526caa4df80e, 0x15ed2665f6bd3375, 0xece71fe59752dcdb, 0x9a55af060e7a37a,
	0x1e76e152a404cf7b, 0x2736de88a2bb485d, 0xd5171928029c46d8, 0xf9883a9a7f4edde7,
	0x458da171c25a05d1, 0xd0f502db40319345, 0x1447e08bfda006f, 0x5b4ab52967ba57a0,
	0x71c79ce53a1b3bee, 0x634f0f860f6a06af, 0x9122322cdd0e38a, 0x5efdc5574ebc5eb2,
	0xaa75c7c8e5ac8c4e, 0x35de30a85c29d519, 0x512c3b0634b99854, 0x660f0d4f370d44f4,
	0xaa07a81256fc8f2c, 0x942e0a917b8b3867, 0x5a243abedf2361b0, 0x9315a550edf3dcd,
	0xf6aa07424e4ffa40, 0xa198de67103b6db0, 0x1f55171fbd582764, 0xf9771c83a121db5e,
	0x8e3a96ea54644932, 0x7bad02672224de0a, 0x56fa7f13df687617, 0x8c31d183155711f5,
	0x3f62525d8ccf8a4e, 0x87be47ca3c00df4, 0x12bab94a6fadfb3a, 0xe6bb75b7bb87c4b8,
	0x27ad62dac633fc57, 0xf45f74285777f093, 0x2105dc6d03b22ba0, 0xafec404cec194584,
	0x7cf1d60ca3114ff5, 0xa4679ea9b6cd3a2a, 0xb9cf71eeca9a3772, 0xdfed363220087ff2,
	0xe3bc29be85fd63c, 0x18ce5d041716080b, 0x267db7585bbc6721, 0xf19d288d4946e936,
	0xa6bcb748ee73768, 0x8f5ee6a0afae1600, 0x6274287ffeefe37e, 0xb231f81a02c870df,
	0xb3b42a8a4f9af9ca, 0x62712aae87a8135e, 0x61c3ebb392ecce2f, 0xf36c5b0703e53153,
	0xe95dae6487af795, 0x908c4d671f3adce7, 0x16c26678c6466e1d, 0x4316a0c54bd9508e,
	0x296079bc556151fd, 0xa8ef54e4494ee19, 0x6b81df07c2c9ee06, 0xd8a29b23c58df45a,
	0x68a0098c9129a742, 0x8211e32acde545b3, 0x8d41a67e03524a84, 0xd3e4dc282417279b,
	0x6434fdd9455242c5, 0x51103f4fa800074c, 0x4dc824517b752a7f, 0x402780ff6427295e,
	0x12cdaf868695a927, 0x5fb2bdad8bc933bd, 0xc239be1aaf1a317f, 0xa024f97529a6cb39,
	0xd7a9a7e52a090783, 0xe2e74b0106b114e7, 0x31a0d2898d896a1d, 0x4b59cc1ddb77931b,
	0xa4193e4b7ca1e56f, 0x85d47266d2b92319, 0xcded56b6ed5edcb4, 0xdb438c5bbc3572e1,
	0x57d9c5a9c49960d7, 0x9ba7aeb933c19f44, 0x598e40c4254e1979, 0xf69542a77900ffbb,
	0x7aecfb4cadbd3f0f, 0xa5ab906ac497a416, 0x9b4234a0fc194613, 0x3a8b07924a1f1c3c,
	0xc6817c4b7b26dbd9, 0x5b5fe8bc04f55ba2, 0x4abd439ba09258be, 0xd4464b3f344f351c,
	0x76d809177f032b53, 0x37e535b7d4c9e18, 0x95ecfbefad98fa1a, 0xb99f3b564a2be33c,
	0xe0f5cb26ac8d58a3, 0x9bef193dab11f550, 0x2cafebaa20483d2c, 0x206287c0237cdf7a,
	0x72e0e8166f114123, 0x6d8e07841e9f0b89, 0x7197b29ec8e44eb5, 0x620abe78d6c83fd3,
	0x41c972da93953ffa, 0x3cb7806f441c54cd, 0x269625916907bb83, 0x71158f8d8e5ea268,
	0xf6eaab8e9580bf04, 0x9cf4f1e2bc1cb70c, 0x9ca50efcf2d1d0c9, 0x26578896c6d569cb,
	0xce6d95f6c7773516, 0x737d34ba64ca648, 0x37003f8fd9f43fe7, 0x2e8a8ce9d821b7b3,
}

const (
	maskS uint64 = 0x0003590703530000
	maskA uint64 = 0x0000d90303530000
	maskL uint64 = 0x0000d90003530000

	MinSize int = (1 << 11) // 2^11 = 2KB
	MaxSize int = (1 << 16) // 2^16 = 64KB
)

// Compute calculates chunk boundary over `buf` and returns index of last byte
// of a chunk.
func Compute(buf []byte) int {
	fp := uint64(0)
	i := MinSize
	n := len(buf)
	normalSize := int(1 << 13) // 2^13 = 8KB

	if n <= MinSize {
		return n
	}

	if n >= MaxSize {
		n = MaxSize
	} else if n <= normalSize {
		normalSize = n
	}

	for ; i < normalSize; i++ {
		fp = (fp << 1) + gear[buf[i]]
		if (fp & maskS) == 0 {
			return i
		}
	}

	for ; i < n; i++ {
		fp = (fp << 1) + gear[buf[i]]
		if (fp & maskL) == 0 {
			return i
		}
	}

	return i
}
