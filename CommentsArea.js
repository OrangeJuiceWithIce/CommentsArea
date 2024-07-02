var currentpage=1;
var comments=document.getElementsByClassName("comment");
const PrevPage=document.getElementById("prev-page");
const NextPage=document.getElementById("next-page");
const deletebuttons=document.getElementsByClassName("delete-button");
//给每个delete按钮添加事件
for(let i=0;i<deletebuttons.length;i++){
    deletebuttons[i].addEventListener("click",function(){
        const confirmation = confirm('你确定要删除此评论吗？'); // 弹出确认对话框
        if(confirmation){
            fetch("http://localhost:8080/comment/delete/",{
                method:"POST",
                headers:{
                    "Content-Type":"application/json"
                },
                body:JSON.stringify({
            })
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

ShowComments(currentpage,comments);//默认显示第一页评论

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
    .then(data => {
        console.log("评论成功",data);
        insertcomment(data);
    })
    .catch(error => {
        console.log(error);
    })
})

function insertcomment(data){
    const showarea=document.getElementById("show-field");
    //创建新的评论元素
    var newcomment=document.createElement("div");
    var newusername=document.createElement("span");
    var newcontent=document.createElement("p");
    var newdeletebutton=document.createElement("button");
    //给新元素设置属性
    newcomment.className="comment";
    newusername.className="show-UserName";
    newcontent.className="show-CommentContent";
    newdeletebutton.className="delete-button";
    newdeletebutton.type="button";
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