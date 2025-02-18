package stat

type COMPONENT_TYPE int32

const (
	COMP_BODY COMPONENT_TYPE = iota
	COMP_BRAIN
	COMP_PROPULSION
	COMP_REPAIRUNIT
	COMP_ECM
	COMP_SENSOR
	COMP_CONSTRUCT
	COMP_WEAPON
	COMP_NUMCOMPONENTS
)

const (
	STAT_BODY       = 0x010000
	STAT_BRAIN      = 0x020000
	STAT_PROPULSION = 0x040000
	STAT_SENSOR     = 0x050000
	STAT_ECM        = 0x060000
	STAT_REPAIR     = 0x080000
	STAT_WEAPON     = 0x0a0000
	STAT_RESEARCH   = 0x0b0000
	STAT_TEMPLATE   = 0x0c0000
	STAT_STRUCTURE  = 0x0d0000
	STAT_FUNCTION   = 0x0e0000
	STAT_CONSTRUCT  = 0x0f0000
	STAT_FEATURE    = 0x100000
	STAT_MASK       = 0xffff0000
)

func RefToComponent(ref uint32) (COMPONENT_TYPE, int) {
	switch ref & STAT_MASK {
	case STAT_BODY:
		return COMP_BODY, int(int64(ref) - int64(STAT_BODY))
	case STAT_BRAIN:
		return COMP_BRAIN, int(int64(ref) - int64(STAT_BRAIN))
	case STAT_PROPULSION:
		return COMP_PROPULSION, int(int64(ref) - int64(STAT_PROPULSION))
	case STAT_SENSOR:
		return COMP_SENSOR, int(int64(ref) - int64(STAT_SENSOR))
	case STAT_ECM:
		return COMP_ECM, int(int64(ref) - int64(STAT_ECM))
	case STAT_REPAIR:
		return COMP_REPAIRUNIT, int(int64(ref) - int64(STAT_REPAIR))
	case STAT_WEAPON:
		return COMP_WEAPON, int(int64(ref) - int64(STAT_WEAPON))
	case STAT_RESEARCH:
		return -1, int(int64(ref) - int64(STAT_RESEARCH))
	case STAT_TEMPLATE:
		return -1, int(int64(ref) - int64(STAT_TEMPLATE))
	case STAT_STRUCTURE:
		return -1, int(int64(ref) - int64(STAT_STRUCTURE))
	case STAT_FUNCTION:
		return -1, int(int64(ref) - int64(STAT_FUNCTION))
	case STAT_CONSTRUCT:
		return COMP_CONSTRUCT, int(int64(ref) - int64(STAT_CONSTRUCT))
	case STAT_FEATURE:
		return -1, int(int64(ref) - int64(STAT_FEATURE))
	default:
		return COMP_NUMCOMPONENTS, -1
	}
}
