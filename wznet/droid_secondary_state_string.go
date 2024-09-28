// Code generated by "stringer --type DROID_SECONDARY_STATE"; DO NOT EDIT.

package wznet

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DSS_NONE-0]
	_ = x[DSS_ARANGE_SHORT-1]
	_ = x[DSS_ARANGE_LONG-2]
	_ = x[DSS_ARANGE_OPTIMUM-3]
	_ = x[DSS_REPLEV_LOW-4]
	_ = x[DSS_REPLEV_HIGH-8]
	_ = x[DSS_REPLEV_NEVER-12]
	_ = x[DSS_ALEV_ALWAYS-16]
	_ = x[DSS_ALEV_ATTACKED-32]
	_ = x[DSS_ALEV_NEVER-48]
	_ = x[DSS_HALT_HOLD-64]
	_ = x[DSS_HALT_GUARD-128]
	_ = x[DSS_HALT_PURSUE-192]
	_ = x[DSS_RECYCLE_SET-256]
	_ = x[DSS_ASSPROD_START-512]
	_ = x[DSS_ACCREP_SET-1024]
	_ = x[DSS_ASSPROD_MID-8192]
	_ = x[DSS_ASSPROD_END-262144]
	_ = x[DSS_RTL_REPAIR-524288]
	_ = x[DSS_RTL_BASE-1048576]
	_ = x[DSS_RTL_TRANSPORT-2097152]
	_ = x[DSS_PATROL_SET-4194304]
	_ = x[DSS_CIRCLE_SET-4194560]
	_ = x[DSS_FIREDES_SET-8388608]
}

const _DROID_SECONDARY_STATE_name = "DSS_NONEDSS_ARANGE_SHORTDSS_ARANGE_LONGDSS_ARANGE_OPTIMUMDSS_REPLEV_LOWDSS_REPLEV_HIGHDSS_REPLEV_NEVERDSS_ALEV_ALWAYSDSS_ALEV_ATTACKEDDSS_ALEV_NEVERDSS_HALT_HOLDDSS_HALT_GUARDDSS_HALT_PURSUEDSS_RECYCLE_SETDSS_ASSPROD_STARTDSS_ACCREP_SETDSS_ASSPROD_MIDDSS_ASSPROD_ENDDSS_RTL_REPAIRDSS_RTL_BASEDSS_RTL_TRANSPORTDSS_PATROL_SETDSS_CIRCLE_SETDSS_FIREDES_SET"

var _DROID_SECONDARY_STATE_map = map[DROID_SECONDARY_STATE]string{
	0:       _DROID_SECONDARY_STATE_name[0:8],
	1:       _DROID_SECONDARY_STATE_name[8:24],
	2:       _DROID_SECONDARY_STATE_name[24:39],
	3:       _DROID_SECONDARY_STATE_name[39:57],
	4:       _DROID_SECONDARY_STATE_name[57:71],
	8:       _DROID_SECONDARY_STATE_name[71:86],
	12:      _DROID_SECONDARY_STATE_name[86:102],
	16:      _DROID_SECONDARY_STATE_name[102:117],
	32:      _DROID_SECONDARY_STATE_name[117:134],
	48:      _DROID_SECONDARY_STATE_name[134:148],
	64:      _DROID_SECONDARY_STATE_name[148:161],
	128:     _DROID_SECONDARY_STATE_name[161:175],
	192:     _DROID_SECONDARY_STATE_name[175:190],
	256:     _DROID_SECONDARY_STATE_name[190:205],
	512:     _DROID_SECONDARY_STATE_name[205:222],
	1024:    _DROID_SECONDARY_STATE_name[222:236],
	8192:    _DROID_SECONDARY_STATE_name[236:251],
	262144:  _DROID_SECONDARY_STATE_name[251:266],
	524288:  _DROID_SECONDARY_STATE_name[266:280],
	1048576: _DROID_SECONDARY_STATE_name[280:292],
	2097152: _DROID_SECONDARY_STATE_name[292:309],
	4194304: _DROID_SECONDARY_STATE_name[309:323],
	4194560: _DROID_SECONDARY_STATE_name[323:337],
	8388608: _DROID_SECONDARY_STATE_name[337:352],
}

func (i DROID_SECONDARY_STATE) String() string {
	if str, ok := _DROID_SECONDARY_STATE_map[i]; ok {
		return str
	}
	return "DROID_SECONDARY_STATE(" + strconv.FormatInt(int64(i), 10) + ")"
}