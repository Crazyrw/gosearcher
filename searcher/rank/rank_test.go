package rank

import (
	"fmt"
	"goSearcher/searcher/db"
	"goSearcher/searcher/model"
	"testing"
)

//func TestRank(t *testing.T) {
//	db.ConnectMySql()
//	docRank := Rank("字节", []int{80553124, 81607254, 80564911, 81611175, 81614441, 81617863, 80577687, 81628962, 80595875, 80663301, 80671326, 88403034, 90961846, 90966426, 90966654, 80692135, 90989886, 87104704, 87108463, 87125604, 80391481, 80275918, 80717596, 80725530, 80740312, 84162425, 82354463, 84168889, 82365923, 81169434, 82376059, 82378760, 80308008, 91859369, 80315901, 80320137, 91891398, 85112772, 85124294, 83984238, 85130634, 83989917, 82566205, 92507780, 92544569, 88004477, 85410240, 91152958, 84058794, 84060120, 91383233, 91384276, 89107357, 84082786, 84094531, 89147136, 91100821, 80071275, 80074867, 80082817, 80095642, 81422109, 81444460, 86037947, 91782979, 91789508, 91794982, 92026513, 88826840, 80494328, 80499691, 81901755, 81905114, 91522038, 90256893, 91529813, 90269261, 81213159, 81938249, 81216764, 81218821, 90288837, 81222481, 81223623, 81947260, 81243789, 89552708, 89560709, 89564027, 89904155, 89594997, 87908696, 86301293, 80423390, 80428246, 86321357, 86346188, 86348613, 85876765, 87258586, 87265211, 92592124, 92596610, 87158582, 92406007, 87181958, 91318896, 92430485, 92431744, 91332532, 83938934, 83404380, 83414788, 80101905, 81308601, 88273885, 83158291, 83167052, 88288482, 83173473, 83177369, 88324672, 83185246, 88329147, 90422823, 88462492, 88465361, 88473653, 88476703, 90216389, 90225823, 85583291, 86388003, 85587485, 82243153, 82154751, 86911138, 86914042, 92621214, 92622244, 90115607, 84413678, 84443346, 84446996, 91700753, 86643818, 82946196, 82948825, 89066325, 89074518, 91740177, 89076408, 92069432, 89222627, 92074467, 81840933, 81845011, 81846571, 84970768, 87202810, 84976927, 84978660, 87214393, 84258263, 84264845, 84270508, 84117175, 91963785, 91550714, 84147059, 91565718, 91567669, 91569248, 91997422, 87530068, 87537420, 81673582, 86066408, 86073496, 86087594, 90455733, 91220528, 90472597, 90495084, 80005546, 80014702, 80024154, 89675749, 86424466, 80027779, 90724222, 89685586, 86445046, 88122268, 87011829, 87020185, 88144800, 84501220, 90604624, 84513026, 90619539, 89037419, 89043077, 90638133, 89046118, 90753558, 90757630, 81076774, 81250283, 81285717, 83482526, 87451105, 87457704, 89321001, 87463406, 83497448, 87466783, 87475148, 87476945, 83551717, 87498799, 84705221, 91000978, 81559121, 81567509, 91015533, 81592003, 81597514, 87406023, 91909744, 91920485, 91927362, 87773060, 91939480, 85310273, 82601116, 82606668, 82612695, 88650056, 82648898, 88658506, 88662430, 85611998, 85620477, 88699670, 85648648, 87066604, 87682815, 87686745, 87687837, 81981564, 85158735, 85167798, 91460629, 85178315, 88712348, 82116547, 85186859, 82120446, 88718617, 82123188, 88735829, 88748910, 88754993, 81454915, 89264916, 89292246, 85716130, 81857207, 80814146, 80820554, 81867015, 80822515, 81872821, 92363275, 80829411, 80833195, 92372829, 82406048, 82421884, 80859003, 85651480, 85663612, 85213355, 85235045, 87808805, 87848359, 88620217, 86455302, 86461083, 86465487, 92171215, 82865487, 82893455, 82051844, 86508127, 82064323, 86528073, 82078766, 89418468, 89419235, 89436701, 83660142, 88915745, 83679711, 88935896, 88944321, 88947102, 89385197, 89387849, 89388221, 91409975, 91422889, 81112724, 86732887, 81133096, 82337312, 85018274, 82348895, 85040969, 86120443, 86121929, 87854782, 87856585, 88383823, 86146222, 87875210, 89851682, 83515120, 83536580, 83540732, 89881613, 89888176, 80150731, 80170801, 80178930, 84555295, 84570234, 84572020, 84576851, 80236974, 84581062, 80246855, 80247662, 84356073, 84358604, 84361051, 84388397, 88567352, 88573787, 90830675, 89504205, 89507243, 89540100, 85783431, 85068546, 84014891, 84023373, 84046538, 86771270, 86777669, 86264479, 80541604, 86272332, 86297517, 81376546, 81384181, 81388041, 91836755, 88958098, 83215371, 88982187, 83227968, 86188973, 86196751, 83248008, 85937157, 88179115, 83058240, 83069672, 83082501, 82955610, 91257915, 91269263, 90355760, 91276256, 91276585, 91277011, 89170224, 90366714, 91288339, 89181916, 90377165, 90387928, 85460837, 92233999, 83399463, 88529658, 88532531, 90546576, 92285603, 89712225, 92295183, 89719966, 92299585, 83701194, 89743094, 83721019, 82266187, 82268021, 90158581, 83038064, 81519367, 81525351, 87639846, 81759916, 86855149, 81799310, 86567474, 85359104, 82703385, 82703712, 82726103, 82748201, 83265733, 83277242, 83283186, 90013998, 82042605, 90030973, 87324367, 90050078, 87338327, 89815005, 89816279, 89826550, 89832053, 89766379, 82501794, 89782032, 89790131, 89794222, 82537259, 82545556, 92106541, 92112442, 84800292, 84833711, 84849108, 80959016, 80963777, 80964273, 80973322, 83862545, 83869649, 83894438, 84239243, 84871430, 84242873, 84897468, 81718972, 81720106, 84615872, 81726044, 84624082, 91617881, 84635230, 82456631, 82462198, 80623065, 86651052, 82481757, 86652951, 86657489, 80643383, 86683200, 83616021, 83633366, 83647273, 89962990, 89968558, 84451427, 85809243, 83815955, 85816296, 85500659, 90074012, 83771556, 83774474, 83828429, 91687534, 91687943, 91699707, 85522760, 85538462})
//	fmt.Println(docRank)
//}
func TestCount(t *testing.T) {
	db.ConnectMySql()
	var count int64
	db.MysqlDB.Model(&model.Docs{}).Count(&count)
	fmt.Println(count)
}
