function initActiveMemberHomePage(loggedinusername,firstname,lastname,avatar,membershipstatus,userid,memberid,roles) {
    console.log("Initialising the Active-Member home page for user: "+loggedinusername+", with roles: "+roles);
    if (roles.indexOf("Membership Manager") >= 0) {
        document.getElementById("membershipmanagementmenu").style.visibility = 'visible';
    }
}