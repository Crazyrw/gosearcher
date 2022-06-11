package core

import (
	"fmt"
	"goSearcher/searcher/btree"
	"goSearcher/searcher/utils"
	"log"
	"testing"
	"time"
)

//TODO: Size(node) > blockSize(4k) key="二胎"
func TestCreateInvertIndex(t *testing.T) {
	start := time.Now()
	//1.创建B+树
	tree, err := btree.NewTree("../../searcher/data/index/invert.db")
	if err != nil {
		panic(err)
	}
	defer tree.Close()
	//2.将dictionary1.txt分词结果插入b+树中 使用mmap读取
	_, data := utils.ReadByMMAP("../../searcher/data/terms/dictionary1.txt")
	k, v := utils.GetAllKVS(data)
	fmt.Println(len(k), len(v))
	for i := 0; i < len(k); i++ {
		err := tree.Insert(k[i], v[i])
		if err != nil {
			fmt.Println(k[i])
			fmt.Println(v[i])
			log.Fatalln(err)
		}
	}
	cost := time.Since(start)
	log.Fatalln(cost)
}

/*
	=== RUN   TestCreateMemoryBtree
	12193026 12193026
	1m27.607646013s
	2022/06/06 11:03:40 create memory btree cost  1m27.607646013s
	2022/06/06 11:03:40 query 二胎( 1210 ) cost  18.483µs
	--- PASS: TestCreateMemoryBtree (87.70s)
	PASS
*/
func TestCreateMemoryBtree(t *testing.T) {
	tree := CreateMemoryBtree("../../searcher/data/terms/dictionary1.txt")
	start := time.Now()
	value, ok := tree.Find("二胎")
	cost := time.Since(start)

	if !ok {
		log.Fatalln("无key的信息")
	}
	strs := utils.SplitDocIdsFromValue(fmt.Sprintf("%v", value))
	log.Println("query 二胎(", len(strs), ") cost ", cost)
}

/*
	=== RUN   TestCreateSkipList
	12193026
	2022/06/06 14:14:48 create skiplist cost  19m10.43155128s
	2022/06/06 14:14:57 query 二胎( 10891 ) cost  20.797µs
	--- PASS: TestCreateSkipList (1159.41s)
	PASS
*/
func TestCreateSkipList(t *testing.T) {
	CreateSkipList("../../searcher/data/terms/dictionary1.txt")
}

/*
	=== RUN   TestGetDocuments
	[376.173ms] [rows:1208]
	--- PASS: TestGetDocuments (0.49s)
	PASS
*/
func TestGetDocuments(t *testing.T) {
	//model.ConnectMySql()
	docIds := []int{80565891, 80570435, 80579574, 81631402, 80587333, 81637289, 80592964, 80668735, 90958759, 80685194, 80691127, 80693217, 90978693, 88427614, 90986771, 88435769, 88437551, 88437636, 88440513, 88441774, 87105670, 82812603, 87139001, 82824977, 82837812, 82846939, 80258457, 80268990, 80269018, 80282529, 80293784, 80299051, 80711151, 80713636, 80724378, 80727233, 80742003, 84150851, 84162374, 84169679, 82363039, 82363384, 90551688, 90555857, 82374599, 90559390, 82378173, 90564198, 84195870, 81185613, 90568931, 84199902, 90570858, 81192477, 82391711, 81195899, 82399612, 90587342, 90597110, 83950567, 91877300, 85103911, 85104023, 91895472, 85116919, 85119844, 83975111, 85125604, 85126875, 83988030, 83989507, 85136265, 83991642, 82551203, 85143126, 83999705, 85147019, 85148619, 92502247, 92504008, 82578411, 82584720, 92521822, 92533089, 92540836, 92543453, 92544480, 88000753, 88007371, 88008724, 88009629, 86962530, 88011086, 86965226, 88022559, 88024373, 85404038, 85408538, 85422664, 85435175, 85445238, 91364465, 89459580, 89460166, 89466953, 91372685, 84054972, 89480933, 91391076, 91393339, 89107599, 89493388, 84098077, 89130728, 89142427, 89143146, 89148834, 80901397, 80903605, 80912408, 80919762, 80930045, 91105558, 91107306, 80933273, 80942409, 80087352, 91121611, 81418599, 81424239, 80098879, 91138532, 91139353, 86003459, 91149531, 86020896, 86023944, 86025691, 91763608, 91773143, 91785813, 91792236, 92010978, 91799721, 92017339, 92018517, 92020558, 92033029, 92034680, 88800649, 92035753, 92037914, 88817895, 88819710, 88820644, 88823022, 88828370, 88837571, 88840935, 88842240, 80461594, 80481598, 80489573, 80489926, 80491971, 80494887, 91508430, 91511080, 91515316, 81909708, 91523252, 90257633, 90261107, 81917559, 90261608, 91532732, 90264218, 81200472, 90269278, 81925871, 81928555, 90273314, 81212633, 81214656, 81940453, 90297704, 81242273, 84900610, 84903751, 84905829, 89900886, 89901200, 89577305, 84914053, 89587059, 89914990, 89928718, 87901710, 89931705, 89936958, 87910379, 84946614, 84947714, 87923712, 87932064, 87939645, 87941587, 87949833, 80406282, 80408416, 80409058, 80409174, 80409711, 90909136, 90912727, 86307024, 90917619, 80428634, 80428917, 86313164, 90925263, 87959838, 90927030, 90927652, 90930656, 87967378, 86326962, 80446448, 90943601, 87978355, 91051283, 87990760, 91060743, 91064415, 91065052, 85850227, 85850641, 85851709, 91084861, 91091003, 85860584, 85868740, 85869199, 85869259, 85887829, 90880193, 85279743, 90886966, 90894816, 90897826, 87253396, 87257760, 87260050, 87269585, 87274403, 87275958, 92565674, 92571125, 92583537, 92588332, 92589952, 92596636, 92401971, 92404892, 91302579, 87177317, 87177776, 92417467, 87181487, 92426543, 91321747, 92430154, 83909628, 92432330, 83913701, 83917149, 83921410, 92443709, 83925361, 83404039, 83411243, 83417927, 83421972, 83435460, 80103020, 83441685, 83446480, 80110681, 81302995, 81304332, 80136465, 81327431, 81327909, 80143582, 80143942, 88276436, 88287199, 88294259, 88315850, 88318160, 83176363, 83181123, 92460900, 88324439, 88324492, 88327497, 83187352, 92466749, 92475976, 92483341, 92488677, 90417675, 90420431, 90436728, 85551769, 85554567, 85555184, 86363370, 86366425, 86366731, 90220479, 88493018, 90225946, 88498643, 88499008, 88499110, 86386108, 86388638, 85595151, 86397273, 82211662, 82221741, 82223306, 82229887, 82230950, 82240148, 82156801, 82163644, 82177818, 82189954, 82195356, 86905683, 92604956, 92607558, 90102134, 86916617, 92614047, 92616687, 92618311, 90110977, 90113141, 90113846, 86931012, 90310419, 92641359, 90139987, 90321785, 90330767, 84403623, 84407471, 84415106, 86600936, 82902930, 84433096, 84434173, 84434492, 82924286, 91706544, 86632750, 91709910, 91714652, 91714801, 86639290, 91716441, 89053069, 91719095, 86647959, 91728701, 89066941, 91743322, 89204003, 89205624, 89087077, 89087537, 89088080, 89090939, 92060201, 89216449, 81801218, 92067364, 92068505, 89224465, 81810915, 89233302, 89236403, 92086491, 92087323, 89244808, 89245850, 92094165, 89247954, 92097675, 92098003, 81835234, 81835511, 81836328, 81842719, 84959680, 87208580, 87217794, 84999200, 87230181, 87234354, 87234991, 87236901, 84267241, 87244057, 84106559, 84108482, 84111979, 84119465, 91960072, 84129291, 91967415, 84142072, 87506330, 91552477, 87513524, 91983153, 87514964, 91561333, 91986848, 91566951, 91993903, 87526353, 81653141, 81657187, 81657807, 81659569, 87536076, 91583024, 91583679, 81668697, 81669954, 84304416, 81686321, 81000691, 84323516, 81021406, 81026359, 84345350, 86075117, 86075904, 87560059, 91202547, 86094175, 87578328, 91205221, 90452694, 91205845, 86097139, 87582634, 90459347, 87590779, 87598366, 87599046, 91232316, 90486233, 90486517, 91247788, 90495084, 90495108, 89654378, 89655193, 89657786, 89659194, 89662004, 89668320, 80019203, 80019315, 86419946, 90722994, 86432780, 89687829, 80036311, 89690367, 89690795, 89692164, 86439077, 86439396, 89696327, 86442693, 88102122, 88108482, 88114879, 89642693, 88119172, 89645781, 87006897, 87007037, 88127815, 87027599, 87028498, 87029182, 84502511, 89017263, 90610620, 89019846, 84517531, 90612710, 90612896, 84521598, 90618729, 89028589, 89040833, 90633052, 90634715, 89043115, 90759172, 81063552, 81065415, 81077267, 90780983, 81085654, 90792441, 90798437, 81256398, 81267300, 81272646, 81276914, 81284591, 81287288, 81288245, 83453432, 81292096, 83464625, 83469673, 83469708, 89303373, 89304128, 89307551, 83483010, 87453548, 83494706, 83498666, 87470050, 87480387, 87480455, 87488509, 87497758, 83566416, 83577376, 83586534, 87361001, 83595862, 87364670, 84711449, 87365534, 87367990, 87369213, 87375136, 81550566, 87378864, 84725621, 81557327, 91004595, 81564536, 84737973, 84739596, 81578840, 81582028, 81589887, 91035187, 81593517, 81594854, 91048898, 87401692, 87428243, 87430130, 91912748, 87435703, 91916053, 87758499, 87761491, 87762158, 87440271, 91923746, 87767387, 87771955, 87776616, 91948988, 87793371, 85300400, 87795186, 85326866, 85331354, 82611423, 85348034, 82632577, 82642331, 88656607, 88668426, 85603048, 85608147, 88680643, 85621121, 85624592, 88696414, 88698232, 85633713, 85634829, 87652894, 85646657, 81951689, 87667720, 81962405, 87669865, 87670554, 87685851, 87089181, 87691134, 87093481, 87695551, 81993404, 91458350, 91459174, 85167677, 85168541, 85171425, 82108286, 82116411, 91480594, 88723091, 88727357, 88730249, 88731337, 82145396, 82147971, 88752805, 85958063, 88766169, 85965766, 88774678, 81454561, 85972626, 85974494, 85974995, 89264911, 88786566, 88786706, 88788046, 85983227, 85986298, 89273637, 85987885, 86221931, 85993423, 89289305, 86239721, 89294146, 86242407, 81492223, 81494111, 81499434, 85702927, 85707688, 85708215, 85712645, 85717669, 85724402, 85734624, 80803483, 80804497, 85740105, 80812394, 85747184, 85749356, 92352229, 92353056, 80819220, 92356887, 80822862, 92359868, 80827593, 92363763, 80829130, 80829415, 81877152, 80831924, 80843535, 81892585, 82408588, 82414925, 82416622, 82422042, 82439500, 80861876, 80869932, 80873753, 90651865, 90657605, 90663764, 80888717, 90668240, 80893358, 90669418, 80895077, 80896433, 85672428, 85678634, 85200358, 85692681, 85696806, 85213354, 85217550, 85230528, 85234579, 85236314, 87808084, 85236992, 87823542, 87826874, 87828476, 88052531, 87834337, 88056999, 87836958, 87849770, 88078072, 88082884, 88608087, 88093806, 88094871, 88616623, 88621078, 88621262, 88622049, 88631247, 86460740, 88643833, 88645132, 86471683, 92151185, 92152385, 92153492, 86499341, 82856891, 92167237, 92168750, 82866537, 92178779, 92186691, 82888939, 82889505, 82890268, 82890805, 82051219, 86504855, 82058353, 86510856, 86511432, 86521258, 87711799, 86528963, 82079282, 82080820, 82082469, 82082937, 86536115, 86539439, 82091531, 87733435, 89405957, 87738106, 82651376, 89409126, 82654460, 89413289, 82667221, 89425846, 82669957, 89430233, 83653969, 82685106, 83666189, 82688211, 83675676, 88917170, 83681462, 88927329, 83687844, 83699438, 89354922, 92334029, 92337079, 89394427, 92349282, 92349928, 86813730, 86818644, 91435211, 91437657, 86830291, 86843064, 86849968, 86715394, 81106948, 81106957, 86726162, 86726283, 86727280, 86730226, 86740587, 86740804, 82322438, 82324796, 81138484, 82328405, 82329082, 85013880, 82337468, 82338098, 85019908, 82348737, 82349954, 85034112, 85035103, 85039592, 85040859, 85042170, 85045936, 85049275, 88356053, 86125411, 86131834, 88372793, 86148244, 88386689, 86149371, 83503312, 83507650, 88394760, 87882622, 87882756, 87884349, 88398852, 83516328, 83521055, 89862078, 83535017, 89894651, 82756181, 82756574, 82762670, 82787723, 82790027, 82799688, 80173983, 80177593, 80209060, 84554586, 84575672, 84577675, 80235107, 84590760, 88550738, 88551570, 88556456, 88850924, 88566236, 84393937, 84394164, 88859809, 84395054, 88574144, 88581817, 88581928, 90801823, 90812855, 90813115, 90837945, 90839990, 84660007, 90843654, 84665212, 84670325, 84678355, 89519400, 89520802, 89530119, 89541380, 89545279, 89545942, 85764957, 85766945, 85772784, 85779286, 85064846, 85783877, 85785941, 85070381, 85073012, 85794495, 85796559, 84002295, 84016907, 84022422, 84027421, 84038492, 84039323, 84039408, 86762329, 86766627, 80521710, 86772984, 86255510, 86777693, 86259435, 86787828, 86790433, 80541646, 86799479, 86286148, 81355631, 86289021, 81358945, 81360231, 81365768, 91811749, 91845260, 88956819, 88965841, 86174391, 86174667, 86176828, 88978310, 83222366, 88981317, 83228130, 86188531, 83229909, 85911700, 83232322, 83234664, 88996586, 83237999, 85917980, 83242924, 83248832, 83249324, 88162624, 83249960, 83250184, 85930147, 85945888, 85946694, 88183167, 88189850, 88198832, 83050672, 83051973, 83053436, 83053450, 83067802, 83072838, 83093696, 83094268, 83094344, 91251302, 82965778, 91262229, 91269532, 89163032, 89164418, 82982695, 91273558, 82991280, 82991747, 90365131, 90367556, 90372270, 90373420, 89186043, 89189860, 89190160, 90379798, 90379819, 90380064, 91298293, 89199586, 85455198, 90393804, 90395311, 92202824, 85464173, 85464345, 90398931, 85482614, 85482623, 92234321, 85495682, 92242407, 92243301, 83352984, 83358511, 83361273, 83365624, 83384754, 83391407, 83395869, 80752544, 80758018, 80762220, 80765593, 80765776, 80766517, 88501192, 80771978, 88507467, 88507998, 88512177, 88516623, 80789328, 90517137, 88521846, 80794166, 90522145, 88526296, 80796725, 90527178, 88535351, 90532262, 92281007, 92286468, 92293548, 92299961, 89728876, 83714806, 89745674, 83726702, 83744052, 82260446, 82271786, 82281508, 82284023, 83002011, 90158339, 90162815, 83021753, 90192324, 83039754, 90198814, 90198936, 83049415, 87602656, 81504875, 87617915, 81509992, 87623977, 87633269, 87649410, 81545083, 81781955, 81784444, 81787985, 86860381, 86860999, 86875276, 85363777, 86880173, 86584177, 85368920, 85373657, 85373791, 86888543, 82713756, 86594233, 82716210, 85386149, 82728235, 82738587, 83250184, 83263073, 83272368, 83279432, 83290659, 83291418, 83291699, 83297047, 83298338, 90003575, 82020995, 82030910, 83300550, 90019516, 82036533, 82038693, 82039931, 84776569, 83323715, 87325436, 87326197, 83335375, 87336180, 83339075, 83346599, 89760250, 89763067, 89840388, 89774283, 82510354, 89778513, 89848748, 82517571, 82524651, 82526343, 82532854, 82549279, 92123662, 92129214, 84806905, 84806944, 92132059, 92146126, 92146855, 84828846, 84829026, 84834317, 84845368, 84848736, 80954321, 80957171, 83856079, 80999823, 83863387, 83883140, 83895275, 83898461, 84242366, 84242762, 84244170, 84883585, 81701926, 84889613, 84895477, 84606026, 81713705, 81717369, 84613425, 84618285, 84624402, 81732111, 81734124, 81740806, 91636500, 83119119, 83122845, 83127224, 83130109, 82464597, 83136694, 80614290, 88226865, 88239447, 88245972, 88247238, 86692704, 86693902, 83617065, 86699011, 83638505, 84453086, 84454129, 84458752, 84459787, 89984600, 89989454, 84470092, 89994826, 84476754, 84485561, 84491308, 83800390, 90052208, 91663370, 90056299, 85803472, 85806431, 90060445, 83759502, 90065326, 83762587, 83816131, 91675355, 85819494, 85503141, 83771689, 91683893, 83828098, 83777143, 91690980, 83832870, 91696611, 85838014, 85847762, 85848249, 85848639, 83800390, 85537330, 85537841}
	documents := utils.GetDocuments(docIds)
	fmt.Println(documents)
}
