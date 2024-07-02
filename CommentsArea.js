var currentpage=1;
getcomments(currentpage,-1);
var comments=document.getElementsByClassName("comment");
const PrevPage=document.getElementById("prev-page");
const NextPage=document.getElementById("next-page");
var deletebuttons=document.getElementsByClassName("delete-button");

//默认显示第一页评论
ShowComments(currentpage,comments);

//给每个delete按钮添加事件
for(let i=0;i<deletebuttons.length;i++){
    deletebuttons[i].addEventListener("click",function(){
        const confirmation = confirm('你确定要删除此评论吗？'); // 弹出确认对话框
        if(confirmation){
            deletecomment(this.parentNode.id);
            this.parentNode.remove();
        }
    })
}

//上一页
PrevPage.addEventListener("click",function(){
    if (currentpage>1){
    currentpage--;
    }
    ShowComments(currentpage,comments);
})

//下一页
NextPage.addEventListener("click",function(){
    if(currentpage<comments.length/5){
        currentpage++;
    }
    ShowComments(currentpage,comments);
})

//显示部分评论
function ShowComments(currentpage,comments){
    for(let i=0;i<comments.length;i++){
        comments[i].style.display="none";
    }
    for(let i=currentpage*5-5;i<currentpage*5;i++){
        if(i<comments.length){
            comments[i].style.display="block";
        }
        else{
            break;
        }
    }
}
//添加评论
const submitbutton=document.getElementById("submit");
submitbutton.addEventListener("click",function(event){
    event.preventDefault();
    const name=document.getElementById("UserName").value;
    const content=document.getElementById("CommentContent").value;
    const postdata={
        name:name,
        content:content
    };
    console.log(JSON.stringify(postdata));
    fetch("http://localhost:8080/comment/add",{
        method:"POST",
        headers:{
            "Content-Type":"application/json"
        },
        body:JSON.stringify(postdata)
    })
    .then(response => response.json())
    .then(responsejson => {
        console.log("评论成功",responsejson);
        insertcomment(responsejson.data);
    })
    .catch(error => {
        console.log(error);
    })
})

//获取评论
function getcomments(currentpage,pagesize){
    const url=`http://localhost:8080/comment/get?page=${currentpage}&size=${pagesize}`;
    fetch(url,{
        method:"GET",
        headers:{
            "Content-Type":"application/json"
        }
    })
    .then(response => 
        response.json()
    )
    .then(responsejson => {
        console.log("获取评论成功",responsejson);
        const existingcomments=document.getElementsByClassName("comment");
        for(;existingcomments.length>0;){
            existingcomments[0].remove();
        }
        for(let i=0;i<responsejson.data.total;i++){
            insertcomment(responsejson.data.comments[i]);
        }
        comments=document.getElementsByClassName("comment");
    })
    .catch(error => {
        throw error;
    })
}

//删除评论
function deletecomment(id){
    const url=`http://localhost:8080/comment/delete?id=${id}`
    fetch(url,{
        method:"POST",
        headers:{
            "Content-Type":"application/json"
        }
    })
    .then(response=>
        response.json()
    )
    .then(responsejson=>{
        console.log("删除评论成功",responsejson);
    })
    .catch(error=>{
        throw error;
    })
}

//在页面中插入新评论的元素
function insertcomment(data){
    const showarea=document.getElementById("show-field");
    //创建新的评论元素
    var newcomment=document.createElement("div");
    var newusername=document.createElement("span");
    var newcontent=document.createElement("p");
    var newdeletebutton=document.createElement("button");
    //给新元素设置属性
    newcomment.id=data.id;
    newcomment.className="comment";
    newusername.className="show-UserName";
    newcontent.className="show-CommentContent";
    newdeletebutton.className="delete-button";
    newdeletebutton.type="button";
    //给新按钮添加事件
    newdeletebutton.addEventListener("click",function(){
        const confirmation = confirm('你确定要删除此评论吗？'); // 弹出确认对话框
        if(confirmation){
            deletecomment(this.parentNode.id);
            this.parentNode.remove();
        }
    })
    //给新元素添加文本
    newusername.textContent=data.name;
    newcontent.textContent=data.content;
    newdeletebutton.textContent="删除";
    //排版
    showarea.insertBefore(newcomment,showarea.firstChild);
    newcomment.appendChild(newusername); 
    newcomment.appendChild(newcontent);
    newcomment.appendChild(newdeletebutton);

    comments=document.getElementsByClassName("comment");
    ShowComments(currentpage,comments);
}