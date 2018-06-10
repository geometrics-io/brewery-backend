package brewery

type Brewery struct {
	ID   string
	Name string
}

// statement grabs active memberships and associated names to be parsed and stored in memory for the autocomplete
// select members.id,membership,memstat_id,members.names from member_status inner join members on member_status.id = members.id where active=true;

// StorageClient initiates a StorageService
type Client interface {
	MemberService() MemberStorage
	MembershipService() MembershipStorage
	TransactionService() TransactionStorage
}

// statement grabs active memstat_id and transactions as json
// select member_status.memstat_id,json_agg(json_build_object('timestamp',member_transactions.timestamp,'raw_units',member_transactions.raw_units)) as transactions from member_status join member_transactions on member_transactions.memstat_id = member_status.memstat_id where member_status.active = true group by member_status.memstat_id;
