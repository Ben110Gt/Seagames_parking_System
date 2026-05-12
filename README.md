Motorcycle ticket payment system
Currently, parking at the Seagames dormitory is still using a day-based system, so we decided to build this system

📁 Requirements
1) A login system that separates the roles of the parking attendant and the owner
2) When the motorcycle enters, a ticket will be created and a ticket will be issued (there will be a barcode on the ticket). After the motorcycle leaves, the barcode will be scanned
3) If it exceeds 1 day, there will be a fine of 2000kip per day (for example, in the ticket on 12/4/2026, but when the barcode is scanned on 15/4/2026, there will be a fine of 6000kip)

4) There is a monthly membership system of 60,000kip
5) View daily, weekly, and monthly income
6) Can search for customer motorcycle tickets (in case the customer loses the ticket)

📁 Tech Stack

Language: Go (Golang)

Web Framework: Fiber v2

Database: PostgreSQL

ORM: GORM

APP: flutter

Configuration: Environment-based configuration management.

📁 Project Structure

frontend

backend 
               |-- cmd 
               |      |-- api 
               |             |-- main.go 
               |  
               | 
               |-- internal 
               |           |-- repository 
               |           |-- service 
               |           |-- handler 
               |           |-- routes 
               |           |-- middleware 
               |           |-- server 
               |           |-- domain 
               |           |-- models 
               |           |-- utils
               | 
               |-- database

