resource "gocd_role" "sample" {
  name  = "sample"
  type  = "gocd"
  users = ["nikhil"]
  policy = [{
    "permission" : "allow",
    "action" : "administer",
    "type" : "*",
    "resource" : "*"
  }]
}

resource "gocd_role" "sample_ldap" {
  name           = "sample-ldap"
  type           = "plugin"
  auth_config_id = "ldap-config"
  policy = [{
    "permission" : "allow",
    "action" : "administer",
    "type" : "*",
    "resource" : "*"
  }]
  properties {
    key   = "UserGroupMembershipAttribute"
    value = "testing"
  }
  properties {
    key   = "GroupIdentifiers"
    value = "CN=opts,OU=Groups,OU=TESTCOM,DC=TESTCOM,DC=COM"
  }
  properties {
    key   = "GroupSearchBases"
    value = "OU=Groups,OU=TESTCOM,DC=TESTCOM,DC=COM"
  }
  properties {
    key   = "GroupMembershipFilter"
    value = "(&(member={dn})(cn=opts))"
  }
}


data "gocd_role" "sample" {
  name = gocd_role.sample.id
}

data "gocd_role" "sample_ldap" {
  name = gocd_role.sample_ldap.id
}