# cat

## Usage

	Usage:
	  kafcat cat <topic>[:partition] [flags]

	Flags:
	  -f, --follow         don't quit on end of log; keep following the topic
	  -s, --since string   only return logs newer than a relative duration like 5s, 2m, or 3h, or shorthands 0 and now

	Global Flags:
	  -b, --broker-list string   brokers (default "localhost:9092")
	  -v, --verbose              be verbose

# topics

## Usage

	Usage:
	  kafcat topics [flags]

	Flags:
	  -a, --all          show all possible info; equals -olwrsp
	  -s, --isr          show broker ids of in-sync replicas; implies --partitions
	  -l, --leader       show leader information for each partition; implies --partitions
	  -o, --offsets      show oldest and newest offsets for each partition; implies --partitions
	  -p, --partitions   show partitions
	  -r, --replicas     show broker ids of replicas; implies --partitions
	  -w, --writable     show writable state (valid leader accepting writes) for each partition; implies --partitions

	Global Flags:
	  -b, --broker-list string   brokers (default "localhost:9092")
	  -v, --verbose              be verbose

