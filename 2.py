# _*_ coding:utf-8 _*_
# m = input('请输入月份：')
# m = int(m)
# if 1 <= m <= 12:
#     print('输入的月份合理')
#     if 3 <= m <= 5:
#         print('%d月份是春季' % m)
#     elif 6 <= m <= 8:
#         print('%d月份是夏季' % m)
#     elif 9 <= m <= 11:
#         print('%d月份是秋季' % m)
#     else:
#         print('%d月份是冬季' % m)
# else:
#     print('输入的月份不存在')

# _*_ coding:utf-8 _*_
# m = input('请输入月份：')
# m = int(m)
# a=(('春'if　m<=3　else '夏') if m<=6 else ('秋'if m<=9 else '冬'))
# print(a)
# _*_ coding:utf-8 _*_



# a = {}
# b = int(input('请输入人数：'))
# i = 1
# while i <= b:
#     name = input('输入第%d个学生名字：' % i)
#     score = int(input('输入第%d个学生成绩：' % i))
#     a[name] = score
#     i+= 1
# print(a)
# l=a
# num=(input('''输入您要的序号：
#     1.成绩查询
#     2.成绩录改
#     3.成绩删除
#     4.输出所有人成绩
#     其他任意键退出\n'''))
# if num==1:
#     name1=input('输入您的名字：')
#     if name1 in l:
#         print('%s的成绩为：%d'%(name1,l[name1]))
#     else:
#         print('查无此人')
# elif num==2:
#     name2=(input('输入您的名字：'))
#     if name2 in l:
#         l[name2]=int(input('输入您的成绩：'))
#         print('%s的成绩为：%d'%(name2,l[name2]))
#     else:
#         print('查无此人')
# elif num==3:
#     name3=input('输入您的名字：')
#     if name3 in l:
#         print(l.pop(name3))
#         print(l)
#     else:
#         print('查无此人')
# elif num==4:
#     print(l)
# else:
#     print('系统退出')



# _*_ coding:utf-8 _*_
grade={'a':75,'b':80,'c':90}
exit=input('欢迎使用成绩查询系统！进入按y，其它任意键退出\n')
while exit=='y':
    menu=['1.录入','2.查询','3.修改','4.学生列表','5.退出']
    for features in menu:
        print(features)
    ord=int(input('输入您想要的操作序号：'))
    if ord==1:
        user=input('输入您要录入的名字：')
        grade[user]=input('输入成绩：')
        print('录入完成，正在返回．\n录入的结果为：%s，%s分'%(user,grade[user]))
    elif ord==2:
        user=input('输入您要查询的学生姓名：')
        if user in grade:
            print('%s的成绩为：%s'%(user,grade[user]))
            exit=input('\n查询完毕，输入y继续查询，输入n退出本系统\n')
        else:
            print('此人，请重新输入')
    elif ord==3:
        user=('输入您要修改的学生的姓名：')
        if user in grade:
            grade[user]=input('输入正确的成绩：')
            exit=input

            
