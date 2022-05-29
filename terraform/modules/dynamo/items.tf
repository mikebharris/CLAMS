resource "aws_dynamodb_table_item" "attendees_golden_data" {
  table_name = aws_dynamodb_table.attendees_table.name
  hash_key   = aws_dynamodb_table.attendees_table.hash_key
  count      = 1

  item = jsonencode({
    "Code" : {
      "S" : "5F7BCD"
    },
    "Name" : {
      "S" : "Maximillian Spillage"
    },
    "Email" : {
      "S" : "max.spillage@somedomain.com"
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
        "ToPay" : {
          "N" : "75"
        },
        "Paid" : {
          "N" : "50"
        },
        "PaidDate" : {
          "S" : "28/05/2022"
        },
        "Due" : {
          "N" : "25"
        }
      }
    },
    "StayingLate" : {
      "S" : "Yes"
    },
    "Arrival": {
      "S": "Wednesday"
    }
    "Nights": {
      "N": "5"
    }
  })
}
