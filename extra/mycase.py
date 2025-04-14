from scrapy import Selector
from curl_cffi import requests
import json
import pandas as pd
import scraper_helper
import datetime
from headers import put_headers,json_headers,html_headers


offices = [
    {
        "id": 9811,
        "name": "Abrams Law Group, P.C."
    },
    {
        "id": 378859,
        "name": "Abrams Law Group - Long Island"
    },
    {
        "id": 433719,\
        "name": "Abrams Law Group - LS"
    }
]

office_attorney = {
    "Abrams Law Group, P.C.":{
        "id":9811,"attorney":[
            "Arvin Babadzhanov",
            "Agata Lewis",
            "Boris Abrams",
            "Melanie Abrams"
    ]},


    "Abrams Law Group - Long Island":{
        "id":378859,"attorney":[
            "Arvin Babadzhanov",
            "Agata Lewis",
            "Boris Abrams",
            "Melanie Abrams",
            "Enat Davidov",
            "Yael Mirzayev",
            "Yamilex Carvajal",
            "Belen Reyes",
            "Danielle Babaev"
        ]},

    "Abrams Law Group - LS":{
        "id":433719,"attorney":[
            "Giselle Rojas",
            "Melanie Abrams"
        ]},


}


all_attorneys = {}
df = pd.read_csv('attorneys.csv',dtype=str)
for i in df.to_dict('records'):
    all_attorneys[i['attorney']] = {}
    all_attorneys[i['attorney']]['id'] = int(i['url'].split('/')[-1])
    all_attorneys[i['attorney']]['title'] = i['title']



def AddStaff(row):
    ts = datetime.datetime.now().timestamp()
    ts = int(ts)
    # case_id = '20806463'
    case_id = row['case_id']
    req = requests.get(f"https://abrams-law-group.mycase.com/court_cases/{case_id}/case_contacts_data.json",headers=json_headers,impersonate="chrome133a")
    print(req.url,req.status_code)
    staffs = {}
    for i in req.json()['staff']:
        staffs[i['full_name']] = i
    for att in office_attorney[row['name']]['attorney']:
        if att not in staffs:
            print('Adding', att)
            req1 = requests.get(f"https://abrams-law-group.mycase.com/court_cases/{case_id}/existing_staff?_={ts}",headers=html_headers,impersonate="chrome133a")
            print(req1.status_code,req1.url)
            resp = Selector(text=req1.text)
            auth_token = resp.xpath('//form[@id="user_link_form"]//input[@name="authenticity_token"]/@value').get()
            if auth_token:
                att_id = all_attorneys[att]['id']
                att_title = all_attorneys[att]['title']
                att_name  = att.replace(' ','+')+ f'+({att_title})'
                print(att_name)
                payload = f'_method=put&authenticity_token={auth_token}&existing_id={att_id}&default_rate=&rate_plan_rate=&lawyer_search_name={att_name}&user_link_share_events=on&rate_type=case&court_cases_user%5Bdefault_rate%5D='
                req = requests.put(f"https://abrams-law-group.mycase.com/court_cases/{case_id}/add_staff_link",headers=put_headers,impersonate="chrome133a",data=payload)
                print(req.url,req.status_code)





def RemoveStaff(row):
    case_id = row['case_id']
    sess = requests.Session(headers=json_headers,impersonate="chrome133a")
    req = sess.get(f"https://abrams-law-group.mycase.com/court_cases/{case_id}/case_contacts_data.json")
    print(req.status_code)

    for i in  req.json()['staff']:
        if i['full_name'] not in office_attorney[row['name']]['attorney']:
            req1 = sess.post(f'https://abrams-law-group.mycase.com/court_cases/{case_id}/remove_link.json',json={"user_delete_contact_id":i['id']})
            if req1.status_code == 200:
                    print(f"removed {i['full_name']} from {row['case_name']} case., {req1.text}")
            else:print(req.status_code,req1.text)


main_frame = pd.read_csv('cases.csv')
for i in main_frame.to_dict('records'):
    AddStaff(i)



# RemoveStaff()



# for office in offices:
#     payload = {"firm_user_id":"all","lead_lawyer_id":"all","originating_lawyer_id":"all","practice_area_id":"all",
#         "case_stage_id":"all","office_id":office["id"],"case_status":"open","group_cases_by":"none",
#         "page_limit":1000,"page_number":1,"sort_by":None,"sort_asc":True,
#         "use_date_range":False,"custom_filters":"[]"}

#     req = requests.post("https://abrams-law-group.mycase.com/reporting/case_list_report",headers=json_headers,impersonate="chrome133a",
#         json=payload)
#     print(req.status_code)
#     for i in req.json()['report_data']['All']:
#         office['case_number'] = i['case_number']
#         office['case_id'] = i['id']
#         office['case_name'] = i['name']
#         pd.DataFrame([office]).to_csv('cases.csv',index=False,mode='a')
