#!/usr/bin/python3



import sys,json, mysql.connector,re
import urllib3 
from mysql.connector import Error
from datetime import datetime, timedelta, date ##Change add Date
from os import listdir
from os.path import isfile, join


print("Starting...")
## Preparing redirection
orig_stdout = sys.stdout
process = open('/home/emery/ndt/script_ndt_auto_report/Insertion.log.txt', 'a')
sys.stdout = process



############################################################ Functions
#################### Compute Functions

def ToMbps(total,duration):
	s = duration / 1000000 # s
	b = total*8 # bit
	Mbs = (b/s) / 1000000 # Mbit/s 
	return Mbs

def get_from_StartTime(StartTime):
	y = StartTime.split('-')[0]
	mo = StartTime.split('-')[1]
	d = StartTime.split('-')[2].split('T')[0]
	h = int(StartTime.split('T')[1].split(':')[0])
	mi = int(StartTime.split('T')[1].split(':')[1])
	return (mi,h,d,mo,y)

def nameForm(name):
	return ("Avg"+name,"Min"+name,"Max"+name,"Median"+name)

def  MinMaxMedian(l):
	minl = l[0]
	maxl = l[0]
	medianl = 0.0
	for every in l:
		if every <= minl:
			minl = every
		if every >= maxl:
			maxl = every
	if len(l) %2 == 0:
		ind = (len(l)-1)//2
		medianl = (l[ind]+l[ind+1])/2
	else:
		ind = len(l)//2
		medianl = l[ind]
	return (minl,maxl,medianl)


def Compute(param,infoType,Number_Frame_Test,Frame_Test):
	total = 0
	avgParam = 0
	minParam = None
	maxParam = None
	medParam = 0

	## MinParam and MaxParam
	start = None
	for index in range(Number_Frame_Test):
		if start == None:
			minParam = Frame_Test[index][infoType][param]		
			maxParam = Frame_Test[index][infoType][param]
			start = 1
		if  minParam != None and maxParam != None  and Frame_Test[index][infoType][param] < minParam : ## MinParam
			minParam = Frame_Test[index][infoType][param]
		if minParam != None and maxParam != None  and Frame_Test[index][infoType][param] > maxParam :## MaxParam
			maxParam = Frame_Test[index][infoType][param]						

		total += Frame_Test[index][infoType][param]

	## Avg
	if Number_Frame_Test != 0:
		avgParam = total/Number_Frame_Test
	
	## Median
	if Number_Frame_Test%2 == 0:
		ind = (Number_Frame_Test-1)//2
		medParam = (Frame_Test[ind][infoType][param]+Frame_Test[ind+1][infoType][param])/2
	else:
		ind = Number_Frame_Test // 2
		medParam = Frame_Test[ind+1][infoType][param]

	return (str(avgParam),str(minParam),str(maxParam),str(medParam))


################### Database 

def mysq_version():
    db_Info = connection.get_server_info()
    print("Connected to MySQL Server version ", db_Info)


def Connect():
	connection = mysql.connector.connect(host='localhost',
										database='monitorDB',
										user='root',
										password='Emery@123456789')
	return connection


def insert_into_tcp_bbr_info(table,paramTCP,record,id_test):
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        if id_test == None:
	        	print("Id_test",id_test)
	        	sql_insert_blob_query = """ INSERT INTO """+table+"""
	                           ("""+paramTCP[0]+""","""+paramTCP[1]+""","""+paramTCP[2]+""","""+paramTCP[3]+""") VALUES (%s, %s, %s, %s) """
	        else:
	        	id_test = str(id_test)
		        sql_insert_blob_query = """ Update """+table+"""
	                           set """+paramTCP[0]+"""=%s,"""+paramTCP[1]+"""=%s,"""+paramTCP[2]+"""=%s,"""+paramTCP[3]+"""=%s Where """+table+"""_id="""+id_test+""" """
	        cursor.execute(sql_insert_blob_query,record)
	        connection.commit()
	        print("TCP OR BBR Inserted: ", cursor.rowcount)

	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")

def insert_into_Provider(provider):
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ INSERT INTO Provider
	                           (Provider_ISP,Provider_ORG,Provider_AS_Number,Provider_AS_Name) VALUES (%s, %s, %s ,%s) """
	        cursor.execute(sql_insert_blob_query,provider)
	        connection.commit()
	        print("Provider Inserted: ", cursor.rowcount)

	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")




def insert_into_City(city):
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ INSERT INTO City
	                           (City_Name,City_Region_id) VALUES (%s,%s) """
	        cursor.execute(sql_insert_blob_query,city)
	        connection.commit()
	        print("City Inserted: ", cursor.rowcount)
	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")


def insert_into_Region(region):
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ INSERT INTO Region
	                           (Region_Name,Region_Country_id) VALUES (%s,%s) """
	        cursor.execute(sql_insert_blob_query,region)
	        connection.commit()
	        print("Region Inserted: ", cursor.rowcount)

	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")	


def insert_into_Country(country):
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ INSERT INTO Country
	                           (Country_Name) VALUES ('"""+country+"""') """
	        cursor.execute(sql_insert_blob_query)
	        connection.commit()
	        print("Country Inserted: ", cursor.rowcount)

	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")



def insert_into_Service(service):
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ INSERT INTO Service
	                           (Service_level) VALUES ('"""+service+"""') """
	        cursor.execute(sql_insert_blob_query)
	        connection.commit()
	        print("Service Inserted: ", cursor.rowcount)

	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")

def insert_into_DaySlice(daySlice):
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ INSERT INTO DaySlice
	                           (DaySlice_level) VALUES ('"""+daySlice+"""') """
	        cursor.execute(sql_insert_blob_query)
	        connection.commit()
	        print("DaySlice Inserted: ", cursor.rowcount)

	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")

##Change Test_Year,Test_Month,Test_Day Test_Date
def insert_into_Test(test):
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ INSERT INTO Tests
	                           (Test_UUID,Test_Type,Test_ServerIP,Test_ServerPort,Test_ClientIP,Test_ClientPort,Test_Date,Test_Country_id,Test_Region_id,Test_City_id,Test_Provider_id,Test_BBRInfo_id,Test_TCPInfo_id,Test_Service_id,Test_DaySlice_id) VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) """
	        cursor.execute(sql_insert_blob_query,test)
	        connection.commit()
	        print("Test Inserted: ", cursor.rowcount)

	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")


##Change DayStat_ by DayStat_Date
def insert_into_DayStat(dayStat):
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ INSERT INTO DayStat
	                           (DayStat_Type,DayStat_Date,DayStat_AvgBW,DayStat_MinBW,DayStat_MaxBW,DayStat_MedianBW,DayStat_AvgMinRTT,DayStat_MinMinRTT,DayStat_MaxMinRTT,DayStat_MedianMinRTT,DayStat_Down_Number,DayStat_Up_Number,DayStat_Total) VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) """
	        cursor.execute(sql_insert_blob_query,dayStat)
	        connection.commit()
	        print("DayStat  Inserted: ", cursor.rowcount)

	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")

##Change
def insert_into_inter_DayStat(dayStat):
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ INSERT INTO InterDayStat
	                           (inter_DayStat_Type,inter_DayStat_Date,inter_DayStat_AvgBW,inter_DayStat_MinBW,inter_DayStat_MaxBW,inter_DayStat_MedianBW,inter_DayStat_AvgMinRTT,inter_DayStat_MinMinRTT,inter_DayStat_MaxMinRTT,inter_DayStat_MedianMinRTT,inter_DayStat_Down_Number,inter_DayStat_Up_Number,inter_DayStat_Total) VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) """
	        cursor.execute(sql_insert_blob_query,dayStat)
	        connection.commit()
	        print("InterDayStat  Inserted: ", cursor.rowcount)

	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")


###################  Existence

def isProvider_Exist(asn):
	exist = None
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ SELECT Provider_id FROM Provider WHERE Provider_AS_Number='"""+asn+"""' """
	        cursor.execute(sql_insert_blob_query)
	        exist = cursor.fetchone()  
	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")
	        if exist == None:
	        	return
	        else:
		        return exist[0]



#####################  Fetch data from Database

def get_id_from_table(table1,paramName,param):
	id_need = None
	if param == None:
		return
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ SELECT """+table1+"""_id FROM """+table1+""" WHERE """+paramName+"""='"""+param+"""' """
	       	print(sql_insert_blob_query)
	        cursor.execute(sql_insert_blob_query)
	        id_need = cursor.fetchone()
	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")
	        if id_need == None:
	        	return
	        else:
		        return id_need[0]	


def getLastIndex(table):
	id_need = None
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ SELECT * FROM """+table+"""
	        							ORDER BY """+table+"""_id DESC LIMIT 1 """
	        cursor.execute(sql_insert_blob_query)
	        id_need = cursor.fetchone()      
	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")
	        if id_need == None:
	        	return None
	        else:
	        	return id_need[0]

def get_all(table):
	record = None
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ SELECT * FROM """+table+""" """
	        cursor.execute(sql_insert_blob_query)
	        record = cursor.fetchall()
	        print("Total number of rows in table: ", cursor.rowcount)
	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")
	        return record


def get_all_where(table,t,v):
	val = str(v)
	record = None
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = """ SELECT * FROM """+table+""" where """+t+"""='"""+val+"""' """
	        print(sql_insert_blob_query)
	        cursor.execute(sql_insert_blob_query)
	        record = cursor.fetchone()
	        print("Total number of rows in table: ", cursor.rowcount)
	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")
	        return record

##Change Service flow
def get_service_level(bw):
	b = float(bw) * 8
	w = b / 1000000 
	service = ['0-3','3-5','5-10','10-15','15-25','25-50','50-100','100-1000']
	for levelStr in service:
		l=levelStr.split('-')
		if w > float(l[0]) and w <= float(l[1]):
			return levelStr


def get_level_DaySlice(hour):
	daySlice = ['0-6','6-12','12-18','18-24']
	for h in daySlice:
		l=h.split('-')
		if int(hour) >= int(l[0]) and int(hour) < int(l[1]):
			return h

##Change
def get_Day_Test_Number(dateFormed,loc):
	d = str(dateFormed)
	record = None
	d_num = 0
	u_num = 0
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        if typ == "national":
	        	sql_insert_blob_query = """ SELECT * FROM Tests where Test_Date='"""+d+"""' and Test_Country_id='59' """
	        else:
	        	sql_insert_blob_query = """ SELECT * FROM Tests where Test_Date='"""+d+"""' and Test_Country_id!='59' """
	        
	        cursor.execute(sql_insert_blob_query)
	        record = cursor.fetchall()
	        print("Total number of rows in table: ", cursor.rowcount)
	        for row in record:
	        	if row[2] == 'Download':
	        		d_num += 1
	        	else:
	        		u_num += 1

	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")
	        return (d_num,u_num,d_num+u_num)

##Change 
def get_day_Record(record,t,loc):
	if len(record) == 0:
		return None
	b = []
	r = []
	avgBand = 0.0
	avgRTT = 0.0
	totalBW = 0.0
	totalRTT = 0.0
	for row in record:
		b.append(float(row[1]))
		totalBW += float(row[1])
		r.append(float(row[5]))
		totalRTT += float(row[5])
	avgBand = totalBW/len(record)
	avgRTT = totalRTT/len(record)
	minBand,maxBand,medBand = MinMaxMedian(b)
	minRTT,maxRTT,medRTT = MinMaxMedian(r)

	dateFormed = date(int(sys.argv[2]),int(sys.argv[3]),int(sys.argv[4])).strftime("%Y-%m-%d")
	test = get_Day_Test_Number(dateFormed,loc)
	
	to_save = (t,dateFormed,avgBand,minBand,maxBand,medBand,avgRTT,minRTT,maxRTT,medRTT)+test
	return to_save

##Change
def get_BBRInfo_id_from_Day_Tests(d,t,typ):
	day = str(d)
	record = None
	d_num = 0
	u_num = 0
	try:
	    connection = Connect()
	    if connection.is_connected():
	        cursor = connection.cursor()
	        sql_insert_blob_query = ""
	        if typ == 'national':
	        	sql_insert_blob_query = """ SELECT Test_BBRInfo_id FROM Tests where Test_Date='"""+day+"""' and Test_Type='"""+t+"""' and Test_Country_id='59' """
	        else:
	        	sql_insert_blob_query = """ SELECT Test_BBRInfo_id FROM Tests where Test_Date='"""+day+"""' and Test_Type='"""+t+"""' and Test_Country_id!='59' """
	        
	        print(sql_insert_blob_query)
	        cursor.execute(sql_insert_blob_query)
	        record = cursor.fetchall()
	        print("Total number of rows in table: ", cursor.rowcount)
	except Error as e:
	    print("Error while connecting to MySQL: ", e)
	finally:
	    if connection.is_connected():
	        cursor.close()
	        connection.close()
	        print("MySQL connection is closed")
	        return record


########################  Initail Setup
def setup():
	## Country
	country = [(1, 4, 'AF', 'AFG', 'Afghanistan', 'Afghanistan'),
				(2, 8, 'AL', 'ALB', 'Albania', 'Albanie'),
				(3, 10, 'AQ', 'ATA', 'Antarctica', 'Antarctique'),
				(4, 12, 'DZ', 'DZA', 'Algeria', 'Algérie'),
				(5, 16, 'AS', 'ASM', 'American Samoa', 'Samoa Américaines'),
				(6, 20, 'AD', 'AND', 'Andorra', 'Andorre'),
				(7, 24, 'AO', 'AGO', 'Angola', 'Angola'),
				(8, 28, 'AG', 'ATG', 'Antigua and Barbuda', 'Antigua-et-Barbuda'),
				(9, 31, 'AZ', 'AZE', 'Azerbaijan', 'Azerbaïdjan'),
				(10, 32, 'AR', 'ARG', 'Argentina', 'Argentine'),
				(11, 36, 'AU', 'AUS', 'Australia', 'Australie'),
				(12, 40, 'AT', 'AUT', 'Austria', 'Autriche'),
				(13, 44, 'BS', 'BHS', 'Bahamas', 'Bahamas'),
				(14, 48, 'BH', 'BHR', 'Bahrain', 'Bahreïn'),
				(15, 50, 'BD', 'BGD', 'Bangladesh', 'Bangladesh'),
				(16, 51, 'AM', 'ARM', 'Armenia', 'Arménie'),
				(17, 52, 'BB', 'BRB', 'Barbados', 'Barbade'),
				(18, 56, 'BE', 'BEL', 'Belgium', 'Belgique'),
				(19, 60, 'BM', 'BMU', 'Bermuda', 'Bermudes'),
				(20, 64, 'BT', 'BTN', 'Bhutan', 'Bhoutan'),
				(21, 68, 'BO', 'BOL', 'Bolivia', 'Bolivie'),
				(22, 70, 'BA', 'BIH', 'Bosnia and Herzegovina', 'Bosnie-Herzégovine'),
				(23, 72, 'BW', 'BWA', 'Botswana', 'Botswana'),
				(24, 74, 'BV', 'BVT', 'Bouvet Island', 'Île Bouvet'),
				(25, 76, 'BR', 'BRA', 'Brazil', 'Brésil'),
				(26, 84, 'BZ', 'BLZ', 'Belize', 'Belize'),
				(27, 86, 'IO', 'IOT', 'British Indian Ocean Territory', 'Territoire Britannique de l''Océan Indien'),
				(28, 90, 'SB', 'SLB', 'Solomon Islands', 'Îles Salomon'),
				(29, 92, 'VG', 'VGB', 'British Virgin Islands', 'Îles Vierges Britanniques'),
				(30, 96, 'BN', 'BRN', 'Brunei Darussalam', 'Brunéi Darussalam'),
				(31, 100, 'BG', 'BGR', 'Bulgaria', 'Bulgarie'),
				(32, 104, 'MM', 'MMR', 'Myanmar', 'Myanmar'),
				(33, 108, 'BI', 'BDI', 'Burundi', 'Burundi'),
				(34, 112, 'BY', 'BLR', 'Belarus', 'Bélarus'),
				(35, 116, 'KH', 'KHM', 'Cambodia', 'Cambodge'),
				(36, 120, 'CM', 'CMR', 'Cameroon', 'Cameroun'),
				(37, 124, 'CA', 'CAN', 'Canada', 'Canada'),
				(38, 132, 'CV', 'CPV', 'Cape Verde', 'Cap-vert'),
				(39, 136, 'KY', 'CYM', 'Cayman Islands', 'Îles Caïmanes'),
				(40, 140, 'CF', 'CAF', 'Central African', 'République Centrafricaine'),
				(41, 144, 'LK', 'LKA', 'Sri Lanka', 'Sri Lanka'),
				(42, 148, 'TD', 'TCD', 'Chad', 'Tchad'),
				(43, 152, 'CL', 'CHL', 'Chile', 'Chili'),
				(44, 156, 'CN', 'CHN', 'China', 'Chine'),
				(45, 158, 'TW', 'TWN', 'Taiwan', 'Taïwan'),
				(46, 162, 'CX', 'CXR', 'Christmas Island', 'Île Christmas'),
				(47, 166, 'CC', 'CCK', 'Cocos (Keeling) Islands', 'Îles Cocos (Keeling)'),
				(48, 170, 'CO', 'COL', 'Colombia', 'Colombie'),
				(49, 174, 'KM', 'COM', 'Comoros', 'Comores'),
				(50, 175, 'YT', 'MYT', 'Mayotte', 'Mayotte'),
				(51, 178, 'CG', 'COG', 'Republic of the Congo', 'République du Congo'),
				(52, 180, 'CD', 'COD', 'The Democratic Republic Of The Congo', 'République Démocratique du Congo'),
				(53, 184, 'CK', 'COK', 'Cook Islands', 'Îles Cook'),
				(54, 188, 'CR', 'CRI', 'Costa Rica', 'Costa Rica'),
				(55, 191, 'HR', 'HRV', 'Croatia', 'Croatie'),
				(56, 192, 'CU', 'CUB', 'Cuba', 'Cuba'),
				(57, 196, 'CY', 'CYP', 'Cyprus', 'Chypre'),
				(58, 203, 'CZ', 'CZE', 'Czech Republic', 'République Tchèque'),
				(59, 204, 'BJ', 'BEN', 'Benin', 'Bénin'),
				(60, 208, 'DK', 'DNK', 'Denmark', 'Danemark'),
				(61, 212, 'DM', 'DMA', 'Dominica', 'Dominique'),
				(62, 214, 'DO', 'DOM', 'Dominican Republic', 'République Dominicaine'),
				(63, 218, 'EC', 'ECU', 'Ecuador', 'Équateur'),
				(64, 222, 'SV', 'SLV', 'El Salvador', 'El Salvador'),
				(65, 226, 'GQ', 'GNQ', 'Equatorial Guinea', 'Guinée Équatoriale'),
				(66, 231, 'ET', 'ETH', 'Ethiopia', 'Éthiopie'),
				(67, 232, 'ER', 'ERI', 'Eritrea', 'Érythrée'),
				(68, 233, 'EE', 'EST', 'Estonia', 'Estonie'),
				(69, 234, 'FO', 'FRO', 'Faroe Islands', 'Îles Féroé'),
				(70, 238, 'FK', 'FLK', 'Falkland Islands', 'Îles (malvinas) Falkland'),
				(71, 239, 'GS', 'SGS', 'South Georgia and the South Sandwich Islands', 'Géorgie du Sud et les Îles Sandwich du Sud'),
				(72, 242, 'FJ', 'FJI', 'Fiji', 'Fidji'),
				(73, 246, 'FI', 'FIN', 'Finland', 'Finlande'),
				(74, 248, 'AX', 'ALA', 'Åland Islands', 'Îles Åland'),
				(75, 250, 'FR', 'FRA', 'France', 'France'),
				(76, 254, 'GF', 'GUF', 'French Guiana', 'Guyane Française'),
				(77, 258, 'PF', 'PYF', 'French Polynesia', 'Polynésie Française'),
				(78, 260, 'TF', 'ATF', 'French Southern Territories', 'Terres Australes Françaises'),
				(79, 262, 'DJ', 'DJI', 'Djibouti', 'Djibouti'),
				(80, 266, 'GA', 'GAB', 'Gabon', 'Gabon'),
				(81, 268, 'GE', 'GEO', 'Georgia', 'Géorgie'),
				(82, 270, 'GM', 'GMB', 'Gambia', 'Gambie'),
				(83, 275, 'PS', 'PSE', 'Occupied Palestinian Territory', 'Territoire Palestinien Occupé'),
				(84, 276, 'DE', 'DEU', 'Germany', 'Allemagne'),
				(85, 288, 'GH', 'GHA', 'Ghana', 'Ghana'),
				(86, 292, 'GI', 'GIB', 'Gibraltar', 'Gibraltar'),
				(87, 296, 'KI', 'KIR', 'Kiribati', 'Kiribati'),
				(88, 300, 'GR', 'GRC', 'Greece', 'Grèce'),
				(89, 304, 'GL', 'GRL', 'Greenland', 'Groenland'),
				(90, 308, 'GD', 'GRD', 'Grenada', 'Grenade'),
				(91, 312, 'GP', 'GLP', 'Guadeloupe', 'Guadeloupe'),
				(92, 316, 'GU', 'GUM', 'Guam', 'Guam'),
				(93, 320, 'GT', 'GTM', 'Guatemala', 'Guatemala'),
				(94, 324, 'GN', 'GIN', 'Guinea', 'Guinée'),
				(95, 328, 'GY', 'GUY', 'Guyana', 'Guyana'),
				(96, 332, 'HT', 'HTI', 'Haiti', 'Haïti'),
				(97, 334, 'HM', 'HMD', 'Heard Island and McDonald Islands', 'Îles Heard et Mcdonald'),
				(98, 336, 'VA', 'VAT', 'Vatican City State', 'Saint-Siège (état de la Cité du Vatican)'),
				(99, 340, 'HN', 'HND', 'Honduras', 'Honduras'),
				(100, 344, 'HK', 'HKG', 'Hong Kong', 'Hong-Kong'),
				(101, 348, 'HU', 'HUN', 'Hungary', 'Hongrie'),
				(102, 352, 'IS', 'ISL', 'Iceland', 'Islande'),
				(103, 356, 'IN', 'IND', 'India', 'Inde'),
				(104, 360, 'ID', 'IDN', 'Indonesia', 'Indonésie'),
				(105, 364, 'IR', 'IRN', 'Islamic Republic of Iran', 'République Islamique d''Iran'),
				(106, 368, 'IQ', 'IRQ', 'Iraq', 'Iraq'),
				(107, 372, 'IE', 'IRL', 'Ireland', 'Irlande'),
				(108, 376, 'IL', 'ISR', 'Israel', 'Israël'),
				(109, 380, 'IT', 'ITA', 'Italy', 'Italie'),
				(110, 384, 'CI', 'CIV', 'Côte d''Ivoire', 'Côte d''Ivoire'),
				(111, 388, 'JM', 'JAM', 'Jamaica', 'Jamaïque'),
				(112, 392, 'JP', 'JPN', 'Japan', 'Japon'),
				(113, 398, 'KZ', 'KAZ', 'Kazakhstan', 'Kazakhstan'),
				(114, 400, 'JO', 'JOR', 'Jordan', 'Jordanie'),
				(115, 404, 'KE', 'KEN', 'Kenya', 'Kenya'),
				(116, 408, 'KP', 'PRK', 'Democratic People''s Republic of Korea', 'République Populaire Démocratique de Corée'),
				(117, 410, 'KR', 'KOR', 'Republic of Korea', 'République de Corée'),
				(118, 414, 'KW', 'KWT', 'Kuwait', 'Koweït'),
				(119, 417, 'KG', 'KGZ', 'Kyrgyzstan', 'Kirghizistan'),
				(120, 418, 'LA', 'LAO', 'Lao People''s Democratic Republic', 'République Démocratique Populaire Lao'),
				(121, 422, 'LB', 'LBN', 'Lebanon', 'Liban'),
				(122, 426, 'LS', 'LSO', 'Lesotho', 'Lesotho'),
				(123, 428, 'LV', 'LVA', 'Latvia', 'Lettonie'),
				(124, 430, 'LR', 'LBR', 'Liberia', 'Libéria'),
				(125, 434, 'LY', 'LBY', 'Libyan Arab Jamahiriya', 'Jamahiriya Arabe Libyenne'),
				(126, 438, 'LI', 'LIE', 'Liechtenstein', 'Liechtenstein'),
				(127, 440, 'LT', 'LTU', 'Lithuania', 'Lituanie'),
				(128, 442, 'LU', 'LUX', 'Luxembourg', 'Luxembourg'),
				(129, 446, 'MO', 'MAC', 'Macao', 'Macao'),
				(130, 450, 'MG', 'MDG', 'Madagascar', 'Madagascar'),
				(131, 454, 'MW', 'MWI', 'Malawi', 'Malawi'),
				(132, 458, 'MY', 'MYS', 'Malaysia', 'Malaisie'),
				(133, 462, 'MV', 'MDV', 'Maldives', 'Maldives'),
				(134, 466, 'ML', 'MLI', 'Mali', 'Mali'),
				(135, 470, 'MT', 'MLT', 'Malta', 'Malte'),
				(136, 474, 'MQ', 'MTQ', 'Martinique', 'Martinique'),
				(137, 478, 'MR', 'MRT', 'Mauritania', 'Mauritanie'),
				(138, 480, 'MU', 'MUS', 'Mauritius', 'Maurice'),
				(139, 484, 'MX', 'MEX', 'Mexico', 'Mexique'),
				(140, 492, 'MC', 'MCO', 'Monaco', 'Monaco'),
				(141, 496, 'MN', 'MNG', 'Mongolia', 'Mongolie'),
				(142, 498, 'MD', 'MDA', 'Republic of Moldova', 'République de Moldova'),
				(143, 500, 'MS', 'MSR', 'Montserrat', 'Montserrat'),
				(144, 504, 'MA', 'MAR', 'Morocco', 'Maroc'),
				(145, 508, 'MZ', 'MOZ', 'Mozambique', 'Mozambique'),
				(146, 512, 'OM', 'OMN', 'Oman', 'Oman'),
				(147, 516, 'NA', 'NAM', 'Namibia', 'Namibie'),
				(148, 520, 'NR', 'NRU', 'Nauru', 'Nauru'),
				(149, 524, 'NP', 'NPL', 'Nepal', 'Népal'),
				(150, 528, 'NL', 'NLD', 'Netherlands', 'Pays-Bas'),
				(151, 530, 'AN', 'ANT', 'Netherlands Antilles', 'Antilles Néerlandaises'),
				(152, 533, 'AW', 'ABW', 'Aruba', 'Aruba'),
				(153, 540, 'NC', 'NCL', 'New Caledonia', 'Nouvelle-Calédonie'),
				(154, 548, 'VU', 'VUT', 'Vanuatu', 'Vanuatu'),
				(155, 554, 'NZ', 'NZL', 'New Zealand', 'Nouvelle-Zélande'),
				(156, 558, 'NI', 'NIC', 'Nicaragua', 'Nicaragua'),
				(157, 562, 'NE', 'NER', 'Niger', 'Niger'),
				(158, 566, 'NG', 'NGA', 'Nigeria', 'Nigéria'),
				(159, 570, 'NU', 'NIU', 'Niue', 'Niué'),
				(160, 574, 'NF', 'NFK', 'Norfolk Island', 'Île Norfolk'),
				(161, 578, 'NO', 'NOR', 'Norway', 'Norvège'),
				(162, 580, 'MP', 'MNP', 'Northern Mariana Islands', 'Îles Mariannes du Nord'),
				(163, 581, 'UM', 'UMI', 'United States Minor Outlying Islands', 'Îles Mineures Éloignées des États-Unis'),
				(164, 583, 'FM', 'FSM', 'Federated States of Micronesia', 'États Fédérés de Micronésie'),
				(165, 584, 'MH', 'MHL', 'Marshall Islands', 'Îles Marshall'),
				(166, 585, 'PW', 'PLW', 'Palau', 'Palaos'),
				(167, 586, 'PK', 'PAK', 'Pakistan', 'Pakistan'),
				(168, 591, 'PA', 'PAN', 'Panama', 'Panama'),
				(169, 598, 'PG', 'PNG', 'Papua New Guinea', 'Papouasie-Nouvelle-Guinée'),
				(170, 600, 'PY', 'PRY', 'Paraguay', 'Paraguay'),
				(171, 604, 'PE', 'PER', 'Peru', 'Pérou'),
				(172, 608, 'PH', 'PHL', 'Philippines', 'Philippines'),
				(173, 612, 'PN', 'PCN', 'Pitcairn', 'Pitcairn'),
				(174, 616, 'PL', 'POL', 'Poland', 'Pologne'),
				(175, 620, 'PT', 'PRT', 'Portugal', 'Portugal'),
				(176, 624, 'GW', 'GNB', 'Guinea-Bissau', 'Guinée-Bissau'),
				(177, 626, 'TL', 'TLS', 'Timor-Leste', 'Timor-Leste'),
				(178, 630, 'PR', 'PRI', 'Puerto Rico', 'Porto Rico'),
				(179, 634, 'QA', 'QAT', 'Qatar', 'Qatar'),
				(180, 638, 'RE', 'REU', 'Réunion', 'Réunion'),
				(181, 642, 'RO', 'ROU', 'Romania', 'Roumanie'),
				(182, 643, 'RU', 'RUS', 'Russian Federation', 'Fédération de Russie'),
				(183, 646, 'RW', 'RWA', 'Rwanda', 'Rwanda'),
				(184, 654, 'SH', 'SHN', 'Saint Helena', 'Sainte-Hélène'),
				(185, 659, 'KN', 'KNA', 'Saint Kitts and Nevis', 'Saint-Kitts-et-Nevis'),
				(186, 660, 'AI', 'AIA', 'Anguilla', 'Anguilla'),
				(187, 662, 'LC', 'LCA', 'Saint Lucia', 'Sainte-Lucie'),
				(188, 666, 'PM', 'SPM', 'Saint-Pierre and Miquelon', 'Saint-Pierre-et-Miquelon'),
				(189, 670, 'VC', 'VCT', 'Saint Vincent and the Grenadines', 'Saint-Vincent-et-les Grenadines'),
				(190, 674, 'SM', 'SMR', 'San Marino', 'Saint-Marin'),
				(191, 678, 'ST', 'STP', 'Sao Tome and Principe', 'Sao Tomé-et-Principe'),
				(192, 682, 'SA', 'SAU', 'Saudi Arabia', 'Arabie Saoudite'),
				(193, 686, 'SN', 'SEN', 'Senegal', 'Sénégal'),
				(194, 690, 'SC', 'SYC', 'Seychelles', 'Seychelles'),
				(195, 694, 'SL', 'SLE', 'Sierra Leone', 'Sierra Leone'),
				(196, 702, 'SG', 'SGP', 'Singapore', 'Singapour'),
				(197, 703, 'SK', 'SVK', 'Slovakia', 'Slovaquie'),
				(198, 704, 'VN', 'VNM', 'Vietnam', 'Viet Nam'),
				(199, 705, 'SI', 'SVN', 'Slovenia', 'Slovénie'),
				(200, 706, 'SO', 'SOM', 'Somalia', 'Somalie'),
				(201, 710, 'ZA', 'ZAF', 'South Africa', 'Afrique du Sud'),
				(202, 716, 'ZW', 'ZWE', 'Zimbabwe', 'Zimbabwe'),
				(203, 724, 'ES', 'ESP', 'Spain', 'Espagne'),
				(204, 732, 'EH', 'ESH', 'Western Sahara', 'Sahara Occidental'),
				(205, 736, 'SD', 'SDN', 'Sudan', 'Soudan'),
				(206, 740, 'SR', 'SUR', 'Suriname', 'Suriname'),
				(207, 744, 'SJ', 'SJM', 'Svalbard and Jan Mayen', 'Svalbard etÎle Jan Mayen'),
				(208, 748, 'SZ', 'SWZ', 'Swaziland', 'Swaziland'),
				(209, 752, 'SE', 'SWE', 'Sweden', 'Suède'),
				(210, 756, 'CH', 'CHE', 'Switzerland', 'Suisse'),
				(211, 760, 'SY', 'SYR', 'Syrian Arab Republic', 'République Arabe Syrienne'),
				(212, 762, 'TJ', 'TJK', 'Tajikistan', 'Tadjikistan'),
				(213, 764, 'TH', 'THA', 'Thailand', 'Thaïlande'),
				(214, 768, 'TG', 'TGO', 'Togo', 'Togo'),
				(215, 772, 'TK', 'TKL', 'Tokelau', 'Tokelau'),
				(216, 776, 'TO', 'TON', 'Tonga', 'Tonga'),
				(217, 780, 'TT', 'TTO', 'Trinidad and Tobago', 'Trinité-et-Tobago'),
				(218, 784, 'AE', 'ARE', 'United Arab Emirates', 'Émirats Arabes Unis'),
				(219, 788, 'TN', 'TUN', 'Tunisia', 'Tunisie'),
				(220, 792, 'TR', 'TUR', 'Turkey', 'Turquie'),
				(221, 795, 'TM', 'TKM', 'Turkmenistan', 'Turkménistan'),
				(222, 796, 'TC', 'TCA', 'Turks and Caicos Islands', 'Îles Turks et Caïques'),
				(223, 798, 'TV', 'TUV', 'Tuvalu', 'Tuvalu'),
				(224, 800, 'UG', 'UGA', 'Uganda', 'Ouganda'),
				(225, 804, 'UA', 'UKR', 'Ukraine', 'Ukraine'),
				(226, 807, 'MK', 'MKD', 'The Former Yugoslav Republic of Macedonia', 'L''ex-République Yougoslave de Macédoine'),
				(227, 818, 'EG', 'EGY', 'Egypt', 'Égypte'),
				(228, 826, 'GB', 'GBR', 'United Kingdom', 'Royaume-Uni'),
				(229, 833, 'IM', 'IMN', 'Isle of Man', 'Île de Man'),
				(230, 834, 'TZ', 'TZA', 'United Republic Of Tanzania', 'République-Unie de Tanzanie'),
				(231, 840, 'US', 'USA', 'United States', 'États-Unis'),
				(232, 850, 'VI', 'VIR', 'U.S. Virgin Islands', 'Îles Vierges des États-Unis'),
				(233, 854, 'BF', 'BFA', 'Burkina Faso', 'Burkina Faso'),
				(234, 858, 'UY', 'URY', 'Uruguay', 'Uruguay'),
				(235, 860, 'UZ', 'UZB', 'Uzbekistan', 'Ouzbékistan'),
				(236, 862, 'VE', 'VEN', 'Venezuela', 'Venezuela'),
				(237, 876, 'WF', 'WLF', 'Wallis and Futuna', 'Wallis et Futuna'),
				(238, 882, 'WS', 'WSM', 'Samoa', 'Samoa'),
				(239, 887, 'YE', 'YEM', 'Yemen', 'Yémen'),
				(240, 891, 'CS', 'SCG', 'Serbia and Montenegro', 'Serbie-et-Monténégro'),
				(241, 894, 'ZM', 'ZMB', 'Zambia', 'Zambie')]
	for c in country:
		insert_into_Country(c[4])

	## Region
	region = ['Alibori','Atacora','Atlantique','Borgou','Collines','Couffo','Donga','Littoral','Mono','Ouémé','Plateau','Zou']
	benin_id = get_id_from_table("Country","Country_Name","Benin")
	for r in region:
		insert_into_Region((r,str(benin_id)))


	## City
	city = [['Banikoara','Gogonou','Kandi','Karimama','Malanville','Ségbana'],['Boukoumbé','Cobli','Kérou','Kouandé',
			'Matéri','Natitingou','Tanguiéta','Toukountouna'],['Abomey-Calavi','Allada','Kpomassè','Ouidah','Sô-Ava','Toffo',
			'Tori','Zê'],['Bembèrèkè','Kalalé','N’Dali','Nikki','Parakou','Pèrèrè','Sinendé','Tchaourou'],['Bantè','Dassa-Zoumè','Glazoué',
			'Ouessè','Savalou','Savè'],['Aplahoué','Djakotomey','Dogbo','Klouékanmey','Lalo','Toviklin'],['Bassila','Copargo','Djougou','Ouaké'],
			['Cotonou'],['Athiémé','Bopa','Comè','Grand-Popo','Houéyogbé','Lokossa'],['Adjarra','Adjohoun','Aguégués','Akpro-Missérété','Avrankou',
			'Bonou','Dangbo','Porto-Novo','Sèmè-Podji'],['Ifangni','Adja-ouèrè','Kétou','Pobè','Sakété'],['Abomey','Agbangnizoun','Bohicon',
			'Covè','Djidja','Ouinhi','Zagnanado','Za-kpota','Zogbodomè']]

	for reg_id in range(len(city)):
		r_id = get_id_from_table("Region","Region_Name",region[reg_id])
		for c in city[reg_id]:
			insert_into_City((c,str(r_id)))



	## Service
	service = ['0-3','3-5','5-10','10-15','15-25','25-50','50-100','100-1000']
	for s in service:
		insert_into_Service(s)

	## DaySlice
	daySlice = ['0-6','6-12','12-18','18-24']
	for d in daySlice:
		insert_into_DaySlice(d)




############################################################## End Functions


##############################################################  Data Processing
# Get options {filename,,setup,dayStat}

option = sys.argv[1]

###### Setup Configuration
if option == "setup":
	print("Start Setup...")
	setup()
	print("End Setup...")
	exit()


####### Dayly Stat Configuration

if option == "dayStat":
	print("Start Recording This days's("+sys.argv[4]+") test...")
	today = date(int(sys.argv[2]),int(sys.argv[3]),int(sys.argv[4])).strftime("%Y-%m-%d")
	print(today)
	dayDownTest = get_BBRInfo_id_from_Day_Tests(today,"Download","national")
	dayUpTest = get_BBRInfo_id_from_Day_Tests(today,"Upload","national")
	print(dayDownTest,dayUpTest)
	if dayDownTest == None or dayUpTest == None:
		sys.exit()
	record = []
	for t in dayDownTest:
		s = get_all_where('BBRInfo','BBRInfo_id',int(t[0]))
		record.append(s)
	if record == []:
		sys.exit()
	dayRecord = get_day_Record(record,"Download","national")
	if dayRecord == None:
		sys.exit()
	print("Start Insertion of Daily Record.")
	insert_into_DayStat(dayRecord)
	print("End of Insertion of Daily Record.")
	record = []
	for t in dayUpTest:
		s = get_all_where('BBRInfo','BBRInfo_id',int(t[0]))
		record.append(s)
	if record == []:
		sys.exit()
	dayRecord = get_day_Record(record,"Upload","national")
	if dayRecord == None:
		sys.exit()
	print("Start Insertion of Daily Record.")
	insert_into_DayStat(dayRecord)
	print("End of Insertion of Daily Record.")
	sys.exit()
	print("\n")
	print("\n")
	print("\n")
	print("\n")

	## Intertnational DayStat
	dayDownTest = get_BBRInfo_id_from_Day_Tests(sys.argv[4],"Download","international")
	dayUpTest = get_BBRInfo_id_from_Day_Tests(sys.argv[4],"Upload","international")
	if dayDownTest == None or dayUpTest == None:
		sys.exit()
	record = []
	for t in dayDownTest:
		s = get_all_where('BBRInfo','BBRInfo_id',int(t[0]))
		record.append(s)
	if record == []:
		sys.exit()
	dayRecord = get_day_Record(record,"Download","international")
	if dayRecord == None:
		sys.exit()
	print("Start Insertion of International Daily Record.")
	insert_into_DayStat(dayRecord)
	print("End of Insertion of International Daily Record.")
	record = []
	for t in dayUpTest:
		s = get_all_where('BBRInfo','BBRInfo_id',int(t[0]))
		record.append(s)
	if record == []:
		sys.exit()
	dayRecord = get_day_Record(record,"Upload","international")
	if dayRecord == None:
		sys.exit()
	print("Start Insertion of International Daily Record.")
	insert_into_inter_DayStat(dayRecord)
	print("End of Insertion of International Daily Record.")
	sys.exit()
	print("\n")
	print("\n")
	




####### Test Data processing

######################## 1st Phase  Getting json data 
print("Getting Data from Json File....")
with open(option, "r") as read_file:
   data = json.load(read_file)
TestType = ""
for i in data:
    if i == "Upload":
    	TestType = "Upload"
    else:
    	TestType = "Download"
ServerIP = data['ServerIP']
ServerPort = data['ServerPort']
ClientIP = data['ClientIP']
ClientPort = data['ClientPort']
StartTime = data['StartTime']
EndTime = data['EndTime']
Measure = data[TestType]
UUID = Measure['UUID']
Frame_Test = Measure['ServerMeasurements']
Number_Frame_Test = len(Frame_Test)
ElapsedTime = Frame_Test[Number_Frame_Test-1]['TCPInfo']['ElapsedTime']
TCP_Params = []
TCP_Useless_Params = ['State','CAState','Probes','WScale','AppLimited','Fackets','LastDataSent','LastDataRecv','LastAckRecv','LastAckSent',
'SndSsThresh','SndCwnd','RcvRTT','RcvSpace','PacingRate','MaxPacingRate','MinRTT','SndWnd','ElapsedTime','Options','RcvOooPack']
BBR_Params = ['BW','MinRTT']
#Get TCP Params
for i in Frame_Test[0]['TCPInfo']:
	TCP_Params.append(i)
# Remove useless  Params
for p in TCP_Useless_Params:
	TCP_Params.remove(p)

print("Data getted...")
############################   2nd Phase Process insertion

### Provider process
print("Provider Process-Searching about IPInfo...")
url = 'http://ip-api.com/json/'+ClientIP
print(url)
http = urllib3.PoolManager()
response = http.request('GET',url)
m_bin = response.data
my_ip_data = json.loads(m_bin.decode('utf8'))
print(my_ip_data)
asSplitted = my_ip_data['as'].split()
asnum = asSplitted[0]
asname = my_ip_data['as'][len(asnum)+1:]
asname = asname.replace("'","")
isp = my_ip_data['isp']
isp = isp.replace("'","")
org = my_ip_data['org']
org = org.replace("'","")
if org == '':
	org = isp
print(asnum,asname)

provider_id = isProvider_Exist(asnum)
if provider_id != None:
	print("This provider "+isp+" already exist")
	pass 
else:
	record = (isp,org,asnum,asname)
	print("Insertion of New Provider detected")
	insert_into_Provider(record)



##### Fill  BBRInfo
print("Insertion of BBRInfo into database..")
bbr_data = []
bbr_last_index=getLastIndex('BBRInfo')
if bbr_last_index != None:
	for b in range(len(BBR_Params)):
		normalizeName = nameForm(BBR_Params[b])
		print(normalizeName)
		data = Compute(BBR_Params[b],'BBRInfo',Number_Frame_Test,Frame_Test)
		print(data)
		bbr_data.append(data)
		if b == 0:
			bbr_last_index = None
		insert_into_tcp_bbr_info('BBRInfo',normalizeName,data,bbr_last_index)
		bbr_last_index=getLastIndex('BBRInfo')
else:
	for b in range(len(BBR_Params)):
		normalizeName = nameForm(BBR_Params[b])
		print(normalizeName)
		data = Compute(BBR_Params[b],'BBRInfo',Number_Frame_Test,Frame_Test)
		print(data)
		bbr_data.append(data)
		insert_into_tcp_bbr_info('BBRInfo',normalizeName,data,bbr_last_index)
		bbr_last_index=getLastIndex('BBRInfo')

print("Insertion of BBRInfo Ended")


##### Fill TCPInfo
print("Insertion of TCPInfo into database..")
tcp_last_index=getLastIndex('TCPInfo')
if tcp_last_index != None:
	for b in range(len(TCP_Params)):
		normalizeName = nameForm(TCP_Params[b])
		print(normalizeName)
		data = Compute(TCP_Params[b],'TCPInfo',Number_Frame_Test,Frame_Test)
		print(data)
		if b == 0:
			tcp_last_index = None
		insert_into_tcp_bbr_info('TCPInfo',normalizeName,data,tcp_last_index)
		tcp_last_index=getLastIndex('TCPInfo')
else:
	for b in range(len(TCP_Params)):
		normalizeName = nameForm(TCP_Params[b])
		print(normalizeName)
		data = Compute(TCP_Params[b],'TCPInfo',Number_Frame_Test,Frame_Test)
		print(data)
		insert_into_tcp_bbr_info('TCPInfo',normalizeName,data,tcp_last_index)
		tcp_last_index=getLastIndex('TCPInfo')

print("Insertion of BBRInfo Ended")






### Fill  Tests
print("Getting Location Ids")
country_id = get_id_from_table('Country','Country_Name',my_ip_data['country'])
print(country_id)
if country_id == None:
	print("Country is None: "+my_ip_data['country'])
	sys.exit()
city_id = get_id_from_table('City','City_Name',my_ip_data['city'])
print(city_id)
if city_id == None:
	print("City is None: "+my_ip_data['city'])
	sys.exit()
region_id = get_id_from_table('Region','Region_Name',my_ip_data['regionName'].split()[0])
print(region_id)
if region_id == None:
	print("Region is None: "+my_ip_data['regionName'])
	sys.exit()
provider_id = get_id_from_table('Provider','Provider_ISP',isp) # normalement my_ip_data['isp']
if provider_id == None:
	print("Provider is None..")
	sys.exit()
print(provider_id)

print(bbr_data)
lev = get_service_level(bbr_data[0][0],)
print("lev",lev)
Test_Service_id =  get_id_from_table('Service','Service_level',lev)
print("Test_Service_id:",Test_Service_id)
print("Starting:",StartTime)
dataTest = get_from_StartTime(StartTime)
print("hour ",dataTest[1])
Test_DaySlice = get_level_DaySlice(dataTest[1])
Test_DaySlice_id = get_id_from_table('DaySlice','DaySlice_level',Test_DaySlice)
print('daySlice id ',Test_DaySlice_id)
print(bbr_last_index,tcp_last_index)

##Change
test_date_provisoir = date(int(dataTest[4]),int(dataTest[3]),int(dataTest[2])) ##Change
test_date = test_date_provisoir.strftime("%Y-%m-%d")
print("test_date:",test_date) ##Change
record = (UUID,TestType,ServerIP,ServerPort,ClientIP,ClientPort,test_date,country_id,region_id,city_id,provider_id,bbr_last_index,tcp_last_index,Test_Service_id,Test_DaySlice_id)
##Change End

print(record)
insert_into_Test(record)
print("Test insertion ended..")
print("End of Insertion of Test: "+UUID+" started at: "+StartTime)

######### close log
sys.stdout = orig_stdout
process.close()



