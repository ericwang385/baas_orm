%{
package expr
%}
%union {
  offset int
  node *AstNode
  text string
}

%token NULL
%token INT
%token STR
%token BOOL
%token FLOAT
%token CONST
%token OR
%token AND
%token NOT
%token LIKE
%token NEQ
%token GT
%token LT
%token GTE
%token LTE
%token EQ
%token ADD
%token MINUS
%token MUL
%token DIV
%token CONTAINS
%token ID
%token IND
%token COMMA
%token ANY
%token FUNC
%token LP
%token RP
%token DOLLAR
%left OR
%left AND
%left NOT
%left GT LT GTE LTE EQ NEQ
%left LIKE CONTAINS
%left ADD MINUS
%left MUL DIV
%left PIPE
%right LP
%right AT
%%

input:    e       { yylex.(*Lexer).parseResult=$1.node;};

e:    INT              { $$.node =newAst(ASTValueInt,$1.text,$1.offset); }
    | STR              { $$.node =newAst(ASTValueText,$1.text,$1.offset); }
    | FLOAT            { $$.node =newAst(ASTValueFloat,$1.text,$1.offset); }
    | BOOL             { $$.node =newAst(ASTValueBool,$1.text,$1.offset); }
    | e AND e          { $$.node =newAst(ASTBinaryOP,"and",$2.offset,$1.node,$3.node); }
    | e OR e           { $$.node =newAst(ASTBinaryOP,"or",$2.offset,$1.node,$3.node); }
    | e ADD e          { $$.node =newAst(ASTBinaryOP,"+",$2.offset,$1.node,$3.node); }
    | e MINUS e        { $$.node =newAst(ASTBinaryOP,"-",$2.offset,$1.node,$3.node); }
    | e DIV e          { $$.node =newAst(ASTBinaryOP,"/",$2.offset,$1.node,$3.node); }
    | e MUL e          { $$.node =newAst(ASTBinaryOP,"*",$2.offset,$1.node,$3.node); }
    | NOT e            { $$.node =newAst(ASTUnaryOP,"not",$1.offset,$2.node); }
    | e GT e           { $$.node =newAst(ASTBinaryOP,">",$2.offset,$1.node,$3.node); }
    | e GTE e          { $$.node =newAst(ASTBinaryOP,">=",$2.offset,$1.node,$3.node); }
    | e LT e           { $$.node =newAst(ASTBinaryOP,"<",$2.offset,$1.node,$3.node); }
    | e LTE e          { $$.node =newAst(ASTBinaryOP,"<=",$2.offset,$1.node,$3.node); }
    | e EQ e           { $$.node =newAst(ASTBinaryOP,"=",$2.offset,$1.node,$3.node); }
    | e NEQ e          { $$.node =newAst(ASTBinaryOP,"!=",$2.offset,$1.node,$3.node); }
    | LP e RP          { $$.node =$2.node;}
    | func_call        { $$.node =$1.node;}
    | e LIKE e         { $$.node =newAst(ASTBinaryOP,"like",$2.offset,$1.node,$3.node);}
 //   | negative 	       { $$.node =$1.node;}
    | AT ID        { $$.node =newAst(ASTVariable,$2.text,$1.offset);}
    | ID        { $$.node =newAst(ASTColumn,$1.text,$1.offset);}

//negative : MINUS INT { $$.node =newAst(ASTUnaryOP,"-" + yylex.(*Lexer).Text(),$2.offset); }
//	| MINUS FLOAT { $$.node =newAst(ASTUnaryOP,"-" + yylex.(*Lexer).Text(),$2.offset); }

func_call :     IDD LP e_list RP { $$.node =newAst(ASTFuncCall,$1.node.Value,$1.offset,$3.node.Children...);}
              | IDD LP RP        { $$.node =newAst(ASTFuncCall,$1.node.Value,$1.offset);}
 //             | IDD              { $$.node =newAst(ASTFuncCall,$1.node.Value,$1.offset);}

IDD: ID {$$.node=newAst(ASTFuncCall,yylex.(*Lexer).Text(),$1.offset);};

e_list:   e  {$$.node =newAst(NULL,"",$1.offset,$1.node);}
        | e_list COMMA e  {$$.node=newAst(NULL,"",$3.offset,append($1.node.Children,$3.node)...);};