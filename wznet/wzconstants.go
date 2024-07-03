package wznet

const (
	// Net-related messages.
	NET_MIN_TYPE                 byte = iota + 33 ///< Minimum-1 valid NET_ type, *MUST* be first.
	NET_PING                                      ///< ping players.
	NET_PLAYER_STATS                              ///< player stats
	NET_TEXTMSG                                   ///< A simple text message between machines.
	NET_PLAYERRESPONDING                          ///< computer that sent this is now playing warzone!
	NET_OPTIONS                                   ///< welcome a player to a game.
	NET_KICK                                      ///< kick a player .
	NET_FIREUP                                    ///< campaign game has started, we can go too.. Shortcut message, not to be used in dmatch.
	NET_COLOURREQUEST                             ///< player requests a colour change.
	NET_FACTIONREQUEST                            ///< player requests a colour change.
	NET_AITEXTMSG                                 ///< chat between AIs
	NET_BEACONMSG                                 ///< place beacon
	NET_TEAMREQUEST                               ///< request team membership
	NET_JOIN                                      ///< join a game
	NET_ACCEPTED                                  ///< accepted into game
	NET_PLAYER_INFO                               ///< basic player info
	NET_PLAYER_JOINED                             ///< notice about player joining
	NET_PLAYER_LEAVING                            ///< A player is leaving, (nicely)
	NET_PLAYER_DROPPED                            ///< notice about player dropped / disconnected
	NET_GAME_FLAGS                                ///< game flags
	NET_READY_REQUEST                             ///< player ready to start an mp game
	NET_REJECTED                                  ///< nope, you can't join
	NET_POSITIONREQUEST                           ///< position in GUI player list
	NET_DATA_CHECK                                ///< Data integrity check
	NET_HOST_DROPPED                              ///< Host has dropped
	NET_SEND_TO_PLAYER                            ///< Non-host clients aren't directly connected to each other, so they talk via the host using these messages.
	NET_SHARE_GAME_QUEUE                          ///< Message contains a game message, which should be inserted into a queue.
	NET_FILE_REQUESTED                            ///< Player has requested a file (map/mod/?)
	NET_FILE_CANCELLED                            ///< Player cancelled a file request
	NET_FILE_PAYLOAD                              ///< sending file to the player that needs it
	NET_DEBUG_SYNC                                ///< Synch error messages, so people don't have to use pastebin.
	NET_VOTE                                      ///< player vote
	NET_VOTE_REQUEST                              ///< Setup a vote popup
	NET_SPECTEXTMSG                               ///< chat between spectators
	NET_PLAYERNAME_CHANGEREQUEST                  ///< non-host human player is changing their name.
	NET_PLAYER_SLOTTYPE_REQUEST                   ///< non-host human player is requesting a slot type change, or a host is asking a spectator if they want to play
	NET_PLAYER_SWAP_INDEX                         ///< a host-only message to move a player to another index
	NET_PLAYER_SWAP_INDEX_ACK                     ///< an acknowledgement message from a player whose index is being swapped
	NET_DATA_CHECK2                               ///< Data2 integrity check
	NET_MAX_TYPE                                  ///< Maximum+1 valid NET_ type, *MUST* be last.
)

const (
	// Game-state-related messages, must be processed by all clients at the same game time.
	GAME_MIN_TYPE       byte = iota + 111 ///< Minimum-1 valid GAME_ type, *MUST* be first.
	GAME_DROIDINFO                        ///< update a droid order.
	GAME_STRUCTUREINFO                    ///< Structure state.
	GAME_RESEARCHSTATUS                   ///< research state.
	GAME_TEMPLATE                         ///< a new template
	GAME_TEMPLATEDEST                     ///< remove template
	GAME_ALLIANCE                         ///< alliance data.
	GAME_GIFT                             ///< a luvly gift between players.
	GAME_LASSAT                           ///< lassat firing.
	GAME_GAME_TIME                        ///< Game time. Used for synchronising, so that all messages are executed at the same gameTime on all clients.
	GAME_PLAYER_LEFT                      ///< Player has left or dropped.
	GAME_DROIDDISEMBARK                   ///< droid disembarked from a Transporter
	GAME_SYNC_REQUEST                     ///< Game event generated from scripts that is meant to be synced

	// The following messages are used for debug mode.
	GAME_DEBUG_MODE             ///< Request enable/disable debug mode.
	GAME_DEBUG_ADD_DROID        ///< Add droid.
	GAME_DEBUG_ADD_STRUCTURE    ///< Add structure.
	GAME_DEBUG_ADD_FEATURE      ///< Add feature.
	GAME_DEBUG_REMOVE_DROID     ///< Remove droid.
	GAME_DEBUG_REMOVE_STRUCTURE ///< Remove structure.
	GAME_DEBUG_REMOVE_FEATURE   ///< Remove feature.
	GAME_DEBUG_FINISH_RESEARCH  ///< Research has been completed.

	// End of debug messages.
	GAME_MAX_TYPE ///< Maximum+1 valid GAME_ type, *MUST* be last.

	// The following messages are used for playing back replays.
	REPLAY_ENDED   ///< A special message for signifying the end of the replay
	REPLAY_ENDED_2 = 255
)

var NetMessageType = map[byte]string{
	NET_MIN_TYPE:                 "NET_MIN_TYPE",
	NET_PING:                     "NET_PING",
	NET_PLAYER_STATS:             "NET_PLAYER_STATS",
	NET_TEXTMSG:                  "NET_TEXTMSG",
	NET_PLAYERRESPONDING:         "NET_PLAYERRESPONDING",
	NET_OPTIONS:                  "NET_OPTIONS",
	NET_KICK:                     "NET_KICK",
	NET_FIREUP:                   "NET_FIREUP",
	NET_COLOURREQUEST:            "NET_COLOURREQUEST",
	NET_FACTIONREQUEST:           "NET_FACTIONREQUEST",
	NET_AITEXTMSG:                "NET_AITEXTMSG",
	NET_BEACONMSG:                "NET_BEACONMSG",
	NET_TEAMREQUEST:              "NET_TEAMREQUEST",
	NET_JOIN:                     "NET_JOIN",
	NET_ACCEPTED:                 "NET_ACCEPTED",
	NET_PLAYER_INFO:              "NET_PLAYER_INFO",
	NET_PLAYER_JOINED:            "NET_PLAYER_JOINED",
	NET_PLAYER_LEAVING:           "NET_PLAYER_LEAVING",
	NET_PLAYER_DROPPED:           "NET_PLAYER_DROPPED",
	NET_GAME_FLAGS:               "NET_GAME_FLAGS",
	NET_READY_REQUEST:            "NET_READY_REQUEST",
	NET_REJECTED:                 "NET_REJECTED",
	NET_POSITIONREQUEST:          "NET_POSITIONREQUEST",
	NET_DATA_CHECK:               "NET_DATA_CHECK",
	NET_HOST_DROPPED:             "NET_HOST_DROPPED",
	NET_SEND_TO_PLAYER:           "NET_SEND_TO_PLAYER",
	NET_SHARE_GAME_QUEUE:         "NET_SHARE_GAME_QUEUE",
	NET_FILE_REQUESTED:           "NET_FILE_REQUESTED",
	NET_FILE_CANCELLED:           "NET_FILE_CANCELLED",
	NET_FILE_PAYLOAD:             "NET_FILE_PAYLOAD",
	NET_DEBUG_SYNC:               "NET_DEBUG_SYNC",
	NET_VOTE:                     "NET_VOTE",
	NET_VOTE_REQUEST:             "NET_VOTE_REQUEST",
	NET_SPECTEXTMSG:              "NET_SPECTEXTMSG",
	NET_PLAYERNAME_CHANGEREQUEST: "NET_PLAYERNAME_CHANGEREQUEST",
	NET_PLAYER_SLOTTYPE_REQUEST:  "NET_PLAYER_SLOTTYPE_REQUEST",
	NET_PLAYER_SWAP_INDEX:        "NET_PLAYER_SWAP_INDEX",
	NET_PLAYER_SWAP_INDEX_ACK:    "NET_PLAYER_SWAP_INDEX_ACK",
	NET_DATA_CHECK2:              "NET_DATA_CHECK2",
	NET_MAX_TYPE:                 "NET_MAX_TYPE",
	GAME_MIN_TYPE:                "GAME_MIN_TYPE",
	GAME_DROIDINFO:               "GAME_DROIDINFO",
	GAME_STRUCTUREINFO:           "GAME_STRUCTUREINFO",
	GAME_RESEARCHSTATUS:          "GAME_RESEARCHSTATUS",
	GAME_TEMPLATE:                "GAME_TEMPLATE",
	GAME_TEMPLATEDEST:            "GAME_TEMPLATEDEST",
	GAME_ALLIANCE:                "GAME_ALLIANCE",
	GAME_GIFT:                    "GAME_GIFT",
	GAME_LASSAT:                  "GAME_LASSAT",
	GAME_GAME_TIME:               "GAME_GAME_TIME",
	GAME_PLAYER_LEFT:             "GAME_PLAYER_LEFT",
	GAME_DROIDDISEMBARK:          "GAME_DROIDDISEMBARK",
	GAME_SYNC_REQUEST:            "GAME_SYNC_REQUEST",
	GAME_DEBUG_MODE:              "GAME_DEBUG_MODE",
	GAME_DEBUG_ADD_DROID:         "GAME_DEBUG_ADD_DROID",
	GAME_DEBUG_ADD_STRUCTURE:     "GAME_DEBUG_ADD_STRUCTURE",
	GAME_DEBUG_ADD_FEATURE:       "GAME_DEBUG_ADD_FEATURE",
	GAME_DEBUG_REMOVE_DROID:      "GAME_DEBUG_REMOVE_DROID",
	GAME_DEBUG_REMOVE_STRUCTURE:  "GAME_DEBUG_REMOVE_STRUCTURE",
	GAME_DEBUG_REMOVE_FEATURE:    "GAME_DEBUG_REMOVE_FEATURE",
	GAME_DEBUG_FINISH_RESEARCH:   "GAME_DEBUG_FINISH_RESEARCH",
	GAME_MAX_TYPE:                "GAME_MAX_TYPE",
	REPLAY_ENDED:                 "REPLAY_ENDED",
}

//go:generate stringer --type DroidOrderSybType

type DroidOrderSybType uint32

const (
	DroidOrderSybTypeObj DroidOrderSybType = iota
	DroidOrderSybTypeLoc
	DroidOrderSybTypeSec
)

//go:generate stringer --type DORDER

type DORDER uint32

const (
	DORDER_NONE             DORDER = iota /**< no order set. */
	DORDER_STOP                           /**< stop the current order. */
	DORDER_MOVE                           /**< 2 - move to a location. */
	DORDER_ATTACK                         /**< attack an enemy. */
	DORDER_BUILD                          /**< 4 - build a structure. */
	DORDER_HELPBUILD                      /**< help to build a structure. */
	DORDER_LINEBUILD                      /**< 6 - build a number of structures in a row (walls + bridges). */
	DORDER_DEMOLISH                       /**< demolish a structure. */
	DORDER_REPAIR                         /**< 8 - repair a structure. */
	DORDER_OBSERVE                        /**< keep a target in sensor view. */
	DORDER_FIRESUPPORT                    /**< 10 - attack whatever the linked sensor droid attacks. */
	DORDER_UNUSED_4                       /**< unused */
	DORDER_UNUSED_2                       /**< unused */
	DORDER_RTB                            /**< return to base. */
	DORDER_RTR                            /**< 14 - return to repair at any repair facility*/
	DORDER_UNUSED_5                       /**< unused */
	DORDER_EMBARK                         /**< 16 - board a transporter. */
	DORDER_DISEMBARK                      /**< get off a transporter. */
	DORDER_ATTACKTARGET                   /**< 18 - a suggestion to attack something i.e. the target was chosen because the droid could see it. */
	DORDER_COMMANDERSUPPORT               /**< Assigns droid to the target commander. */
	DORDER_BUILDMODULE                    /**< 20 - build a module (power, research or factory). */
	DORDER_RECYCLE                        /**< return to factory to be recycled. */
	DORDER_TRANSPORTOUT                   /**< 22 - offworld transporter order. */
	DORDER_TRANSPORTIN                    /**< onworld transporter order. */
	DORDER_TRANSPORTRETURN                /**< 24 - transporter return after unloading. */
	DORDER_GUARD                          /**< guard a structure. */
	DORDER_DROIDREPAIR                    /**< 26 - repair a droid. */
	DORDER_RESTORE                        /**< restore resistance points for a structure. */
	DORDER_SCOUT                          /**< 28 - same as move, but stop if an enemy is seen. */
	DORDER_UNUSED_3                       /**< unused */
	DORDER_UNUSED                         /**< unused */
	DORDER_PATROL                         /**< move between two way points. */
	DORDER_REARM                          /**< 32 - order a vtol to rearming pad. */
	DORDER_RECOVER                        /**< pick up an artifact. */
	DORDER_UNUSED_6                       /**< unused */
	DORDER_RTR_SPECIFIED                  /**< return to repair at a specified repair center. */
	DORDER_CIRCLE           DORDER = 40   /**< circles target location and engage. */
	DORDER_HOLD                           /**< hold position until given next order. */
)

const (
	STAT_BODY       = uint32(0x010000)
	STAT_BRAIN      = uint32(0x020000)
	STAT_PROPULSION = uint32(0x040000)
	STAT_SENSOR     = uint32(0x050000)
	STAT_ECM        = uint32(0x060000)
	STAT_REPAIR     = uint32(0x080000)
	STAT_WEAPON     = uint32(0x0a0000)
	STAT_RESEARCH   = uint32(0x0b0000)
	STAT_TEMPLATE   = uint32(0x0c0000)
	STAT_STRUCTURE  = uint32(0x0d0000)
	STAT_FUNCTION   = uint32(0x0e0000)
	STAT_CONSTRUCT  = uint32(0x0f0000)
	STAT_FEATURE    = uint32(0x100000)
	STAT_MASK       = uint32(0xffff0000)
)

const (
	STRUCTUREINFO_MANUFACTURE = iota
	STRUCTUREINFO_CANCELPRODUCTION
	STRUCTUREINFO_HOLDPRODUCTION
	STRUCTUREINFO_RELEASEPRODUCTION
	STRUCTUREINFO_HOLDRESEARCH
	STRUCTUREINFO_RELEASERESEARCH
)

type GIFT_TYPE uint8

const (
	GIFT_RADAR GIFT_TYPE = iota
	GIFT_DROID
	GIFT_RESEARCH
	GIFT_POWER
	GIFT_STRUCTURE
	GIFT_AUTOGAME
)
