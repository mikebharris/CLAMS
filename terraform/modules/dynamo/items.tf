resource "aws_dynamodb_table_item" "attendees_golden_data" {
  table_name = aws_dynamodb_table.attendees_table.name
  hash_key   = aws_dynamodb_table.attendees_table.hash_key
  count      = 1
  item = jsonencode({
    "AuthCode" : {
      "S" : "5F7BCD"
    },
    "Name" : {
      "S" : "Maximillian Spillage"
    },
    "Email" : {
      "S" : "max.spillage@somedomain.com"
    },
    "Telephone" : {
      "S" : "+1-800-BAMSROCKS"
    },
    "NumberOfKids" : {
      "N" : "1"
    },
    "Diet" : {
      "S" : "I only eat onions on Tuesdays and red lentils on Wednesdays"
    },
    "Financials" : {
      "M" : {
        "AmountToPay" : {
          "N" : "75"
        },
        "AmountPaid" : {
          "N" : "50"
        },
        "DatePaid" : {
          "S" : "28/05/2022"
        },
        "AmountDue" : {
          "N" : "25"
        }
      }
    },
    "StayingLate" : {
      "S" : "Yes"
    },
    "ArrivalDay": {
      "S": "Wednesday"
    }
    "NumberOfNights": {
      "N": "5"
    }
  })
}
