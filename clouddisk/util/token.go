package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var (
	// 1000 primes
	primes = []int{
		20011, 20021, 20023, 20029, 20047, 20051, 20063, 20071, 20089, 20101,
		20107, 20113, 20117, 20123, 20129, 20143, 20147, 20149, 20161, 20173,
		20177, 20183, 20201, 20219, 20231, 20233, 20249, 20261, 20269, 20287,
		20297, 20323, 20327, 20333, 20341, 20347, 20353, 20357, 20359, 20369,
		20389, 20393, 20399, 20407, 20411, 20431, 20441, 20443, 20477, 20479,
		20483, 20507, 20509, 20521, 20533, 20543, 20549, 20551, 20563, 20593,
		20599, 20611, 20627, 20639, 20641, 20663, 20681, 20693, 20707, 20717,
		20719, 20731, 20743, 20747, 20749, 20753, 20759, 20771, 20773, 20789,
		20807, 20809, 20849, 20857, 20873, 20879, 20887, 20897, 20899, 20903,
		20921, 20929, 20939, 20947, 20959, 20963, 20981, 20983, 21001, 21011,
		21013, 21017, 21019, 21023, 21031, 21059, 21061, 21067, 21089, 21101,
		21107, 21121, 21139, 21143, 21149, 21157, 21163, 21169, 21179, 21187,
		21191, 21193, 21211, 21221, 21227, 21247, 21269, 21277, 21283, 21313,
		21317, 21319, 21323, 21341, 21347, 21377, 21379, 21383, 21391, 21397,
		21401, 21407, 21419, 21433, 21467, 21481, 21487, 21491, 21493, 21499,
		21503, 21517, 21521, 21523, 21529, 21557, 21559, 21563, 21569, 21577,
		21587, 21589, 21599, 21601, 21611, 21613, 21617, 21647, 21649, 21661,
		21673, 21683, 21701, 21713, 21727, 21737, 21739, 21751, 21757, 21767,
		21773, 21787, 21799, 21803, 21817, 21821, 21839, 21841, 21851, 21859,
		21863, 21871, 21881, 21893, 21911, 21929, 21937, 21943, 21961, 21977,
		21991, 21997, 22003, 22013, 22027, 22031, 22037, 22039, 22051, 22063,
		22067, 22073, 22079, 22091, 22093, 22109, 22111, 22123, 22129, 22133,
		22147, 22153, 22157, 22159, 22171, 22189, 22193, 22229, 22247, 22259,
		22271, 22273, 22277, 22279, 22283, 22291, 22303, 22307, 22343, 22349,
		22367, 22369, 22381, 22391, 22397, 22409, 22433, 22441, 22447, 22453,
		22469, 22481, 22483, 22501, 22511, 22531, 22541, 22543, 22549, 22567,
		22571, 22573, 22613, 22619, 22621, 22637, 22639, 22643, 22651, 22669,
		22679, 22691, 22697, 22699, 22709, 22717, 22721, 22727, 22739, 22741,
		22751, 22769, 22777, 22783, 22787, 22807, 22811, 22817, 22853, 22859,
		22861, 22871, 22877, 22901, 22907, 22921, 22937, 22943, 22961, 22963,
		22973, 22993, 23003, 23011, 23017, 23021, 23027, 23029, 23039, 23041,
		23053, 23057, 23059, 23063, 23071, 23081, 23087, 23099, 23117, 23131,
		23143, 23159, 23167, 23173, 23189, 23197, 23201, 23203, 23209, 23227,
		23251, 23269, 23279, 23291, 23293, 23297, 23311, 23321, 23327, 23333,
		23339, 23357, 23369, 23371, 23399, 23417, 23431, 23447, 23459, 23473,
		23497, 23509, 23531, 23537, 23539, 23549, 23557, 23561, 23563, 23567,
		23581, 23593, 23599, 23603, 23609, 23623, 23627, 23629, 23633, 23663,
		23669, 23671, 23677, 23687, 23689, 23719, 23741, 23743, 23747, 23753,
		23761, 23767, 23773, 23789, 23801, 23813, 23819, 23827, 23831, 23833,
		23857, 23869, 23873, 23879, 23887, 23893, 23899, 23909, 23911, 23917,
		23929, 23957, 23971, 23977, 23981, 23993, 24001, 24007, 24019, 24023,
		24029, 24043, 24049, 24061, 24071, 24077, 24083, 24091, 24097, 24103,
		24107, 24109, 24113, 24121, 24133, 24137, 24151, 24169, 24179, 24181,
		24197, 24203, 24223, 24229, 24239, 24247, 24251, 24281, 24317, 24329,
		24337, 24359, 24371, 24373, 24379, 24391, 24407, 24413, 24419, 24421,
		24439, 24443, 24469, 24473, 24481, 24499, 24509, 24517, 24527, 24533,
		24547, 24551, 24571, 24593, 24611, 24623, 24631, 24659, 24671, 24677,
		24683, 24691, 24697, 24709, 24733, 24749, 24763, 24767, 24781, 24793,
		24799, 24809, 24821, 24841, 24847, 24851, 24859, 24877, 24889, 24907,
		24917, 24919, 24923, 24943, 24953, 24967, 24971, 24977, 24979, 24989,
		25013, 25031, 25033, 25037, 25057, 25073, 25087, 25097, 25111, 25117,
		25121, 25127, 25147, 25153, 25163, 25169, 25171, 25183, 25189, 25219,
		25229, 25237, 25243, 25247, 25253, 25261, 25301, 25303, 25307, 25309,
		25321, 25339, 25343, 25349, 25357, 25367, 25373, 25391, 25409, 25411,
		25423, 25439, 25447, 25453, 25457, 25463, 25469, 25471, 25523, 25537,
		25541, 25561, 25577, 25579, 25583, 25589, 25601, 25603, 25609, 25621,
		25633, 25639, 25643, 25657, 25667, 25673, 25679, 25693, 25703, 25717,
		25733, 25741, 25747, 25759, 25763, 25771, 25793, 25799, 25801, 25819,
		25841, 25847, 25849, 25867, 25873, 25889, 25903, 25913, 25919, 25931,
		25933, 25939, 25943, 25951, 25969, 25981, 25997, 25999, 26003, 26017,
		26021, 26029, 26041, 26053, 26083, 26099, 26107, 26111, 26113, 26119,
		26141, 26153, 26161, 26171, 26177, 26183, 26189, 26203, 26209, 26227,
		26237, 26249, 26251, 26261, 26263, 26267, 26293, 26297, 26309, 26317,
		26321, 26339, 26347, 26357, 26371, 26387, 26393, 26399, 26407, 26417,
		26423, 26431, 26437, 26449, 26459, 26479, 26489, 26497, 26501, 26513,
		26539, 26557, 26561, 26573, 26591, 26597, 26627, 26633, 26641, 26647,
		26669, 26681, 26683, 26687, 26693, 26699, 26701, 26711, 26713, 26717,
		26723, 26729, 26731, 26737, 26759, 26777, 26783, 26801, 26813, 26821,
		26833, 26839, 26849, 26861, 26863, 26879, 26881, 26891, 26893, 26903,
		26921, 26927, 26947, 26951, 26953, 26959, 26981, 26987, 26993, 27011,
		27017, 27031, 27043, 27059, 27061, 27067, 27073, 27077, 27091, 27103,
		27107, 27109, 27127, 27143, 27179, 27191, 27197, 27211, 27239, 27241,
		27253, 27259, 27271, 27277, 27281, 27283, 27299, 27329, 27337, 27361,
		27367, 27397, 27407, 27409, 27427, 27431, 27437, 27449, 27457, 27479,
		27481, 27487, 27509, 27527, 27529, 27539, 27541, 27551, 27581, 27583,
		27611, 27617, 27631, 27647, 27653, 27673, 27689, 27691, 27697, 27701,
		27733, 27737, 27739, 27743, 27749, 27751, 27763, 27767, 27773, 27779,
		27791, 27793, 27799, 27803, 27809, 27817, 27823, 27827, 27847, 27851,
		27883, 27893, 27901, 27917, 27919, 27941, 27943, 27947, 27953, 27961,
		27967, 27983, 27997, 28001, 28019, 28027, 28031, 28051, 28057, 28069,
		28081, 28087, 28097, 28099, 28109, 28111, 28123, 28151, 28163, 28181,
		28183, 28201, 28211, 28219, 28229, 28277, 28279, 28283, 28289, 28297,
		28307, 28309, 28319, 28349, 28351, 28387, 28393, 28403, 28409, 28411,
		28429, 28433, 28439, 28447, 28463, 28477, 28493, 28499, 28513, 28517,
		28537, 28541, 28547, 28549, 28559, 28571, 28573, 28579, 28591, 28597,
		28603, 28607, 28619, 28621, 28627, 28631, 28643, 28649, 28657, 28661,
		28663, 28669, 28687, 28697, 28703, 28711, 28723, 28729, 28751, 28753,
		28759, 28771, 28789, 28793, 28807, 28813, 28817, 28837, 28843, 28859,
		28867, 28871, 28879, 28901, 28909, 28921, 28927, 28933, 28949, 28961,
		28979, 29009, 29017, 29021, 29023, 29027, 29033, 29059, 29063, 29077,
		29101, 29123, 29129, 29131, 29137, 29147, 29153, 29167, 29173, 29179,
		29191, 29201, 29207, 29209, 29221, 29231, 29243, 29251, 29269, 29287,
		29297, 29303, 29311, 29327, 29333, 29339, 29347, 29363, 29383, 29387,
		29389, 29399, 29401, 29411, 29423, 29429, 29437, 29443, 29453, 29473,
		29483, 29501, 29527, 29531, 29537, 29567, 29569, 29573, 29581, 29587,
		29599, 29611, 29629, 29633, 29641, 29663, 29669, 29671, 29683, 29717,
		29723, 29741, 29753, 29759, 29761, 29789, 29803, 29819, 29833, 29837,
		29851, 29863, 29867, 29873, 29879, 29881, 29917, 29921, 29927, 29947,
		29959, 29983, 29989, 30011, 30013, 30029, 30047, 30059, 30071, 30089,
		30091, 30097, 30103, 30109, 30113, 30119, 30133, 30137, 30139, 30161,
	}
)

type TokenGenerater struct {
	n, e, d int
}

func NewTokenGenerater() *TokenGenerater {
	rand.Seed(int64(time.Now().Second()))

	p, q := primes[int(rand.Int31())%len(primes)], primes[int(rand.Int31())%len(primes)]
	n := p * q
	e := 0x10001
	d := Inverse(e, (p-1)*(q-1))
	return &TokenGenerater{
		n: n,
		e: e,
		d: d,
	}
}

func (tg *TokenGenerater) UID2Token(uid int) string {
	return fmt.Sprintf("%x", Pow(uid+0x10001, tg.e, tg.n))
}

func (tg *TokenGenerater) Token2UID(token string) (int, error) {
	num, err := strconv.ParseInt(token, 16, 0)
	if err != nil {
		return 0, err
	}
	return Pow(int(num), tg.d, tg.n) - 0x10001, nil
}

func Inverse(u, v int) int {
	u3, v3 := u, v
	u1, v1 := 1, 0
	for v3 > 0 {
		q, _ := divmod(u3, v3)
		u1, v1 = v1, u1-v1*q
		u3, v3 = v3, u3-v3*q
	}
	for u1 < 0 {
		u1 += v
	}
	return u1
}

func Pow(e, n, m int) int {
	res := 1
	k, num := 0, n
	for num > 0 {
		num >>= 1
		k++
	}
	for i := 0; i < k; i++ {
		if (n>>uint(i))&1 > 0 {
			res = (res * e) % m
		}
		e = (e * e) % m
	}
	return res
}

func divmod(u, v int) (int, int) {
	return int(u / v), u % v
}
