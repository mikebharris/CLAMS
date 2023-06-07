variable environment {}
variable region {}
variable account_number {}
variable contact {}
variable product {}
variable orchestration {
  default = "https://github.com/mikebharris/clams"
}
variable distribution_bucket {}
variable input_queue_name {
  default = "attendee-input-queue"
}
variable attendees_table_name {
  default = "attendees-datastore"
}
variable certificate_domain{
  default = "events.hacktionlab.org"
}
variable frontend_domain{
  default = "clams.events.hacktionlab.org"
}