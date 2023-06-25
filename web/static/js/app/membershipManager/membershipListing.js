$(document).ready(function () {
    $('#membershiptable').DataTable({
        responsive: true,
        ajax: '/api/listmembers',
        columns: [
            { data: 'FirstName' },
            { data: 'LastName' },
            { data: 'ClubMembershipNumber' },
            { data: 'Email' },
            { data: 'MembershipStatus'},
            {
                data: 'Key',
                render: function (data, type) {
                    let viewURL = '/viewmember?key=' + encodeURIComponent(data);
                    return '<a href="'+viewURL+'">View</a>';
                }
            },
        ],
    });
});
//// UserID	MemberID	Username	AuthUsername	EmailVerified	ClubMembershipNumber	FirstName	LastName	Email	Address	PreferredPhone	EmergencyContact	NAWMembershipNumber	YouthMember	LifeMember	MembershipStatus
