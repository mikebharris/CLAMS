resource "aws_dynamodb_table_item" "attendees_golden_data" {
  table_name = aws_dynamodb_table.attendees_table.name
  hash_key   = aws_dynamodb_table.attendees_table.hash_key
  count      = 1

  item = jsonencode({
    "Code" : {
      "S" : "5F7BCD"
    },
    "Name" : {
      "S" : "Maximillian Schramp"
    },
    "Email" : {
      "S" : "max.schramp@somedomain.com"
    },
    "Phone" : {
      "S" : "+1-800-BAMSROCKS"
    },
    "Kids" : {
      "N" : "1"
    },
    "Diet" : {
      "S" : "I only eat onions on Tuesdays and red lentils on Wednesdays"
    },
    "Financials" : {
      "M" : {
        "To Pay" : {
          "N" : "50"
        },
        "Paid" : {
          "N" : "50"
        },
        "Paid date" : {
          "S" : "28/05/2022"
        }
      }
    },
    "StayingLate" : {
      "S" : "Yes"
    }
  })
}
