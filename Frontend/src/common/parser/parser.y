%lex
%options case-insensitive

%%
\s+             /* skip whitespace */
\n+             /* skip newline */
(true|false)    return 'CONDITION';
(\-)?[0-9]+("."[0-9]+)?\b  return 'NUMBER';
\'(\\.|[^\\'\n])*\' return 'SQM_STRING_LITERAL'
"("             return  '(';
")"             return  ')';
"+"             return  '+';
"-"             return  '-';
"*"             return  '*';
"/"             return  '/';
">="            return  '>=';
"<="            return  '<=';
">"             return  '>';
"<"             return  '<';
"="             return  '=';
"!="            return  '!=';
"and"           return  'AND';
"or"            return  'OR';
"not"           return  'NOT';
"like"          return 'LIKE';


[\$\w\d\u4e00-\u9fa5.]+	{ return 'ID'};
<<EOF>>	return 'EOF';

/lex

/* operator associations and precedence */

%left 'OR'
%left 'AND'
%left 'NOT'
%left 'LIKE'
%left '>' '>=' '<=' '<'   '=' '!='
%left '+' '-'
%left '*' '/'
%right '('

%start end_expression
%%

end_expression
    :   expression EOF
        {return $1;}
    ;

expression
    :   CONDITION
        {$$ = {T:'value',V:($1.toLowerCase() === "true")?"true":"false"};}
    |   NUMBER
        {$$ = {T:'value',V:Number(yytext).toString()};}
    |   SQM_STRING_LITERAL
        {$$ = {T:'value',V:yytext.replaceAll("''","'")};}
    |   ID
        {$$ = yytext.startsWith('$')?{T:'ph',V:yytext}:{T:'col',V:yytext};}
    |   '(' expression ')'
        {$$ = $2;}
    |   func_call
        {$$ = $1;}
    |   expression  '+' expression
        {$$ = {T:'bop',V:'+',C:[$1,$3]};}
    |   expression  '-' expression
        {$$ = {T:'bop',V:'-',C:[$1,$3]};}
    |   expression  '*' expression
        {$$ = {T:'bop',V:'*',C:[$1,$3]};}
    |   expression  '/' expression
        {$$ = {T:'bop',V:'/',C:[$1,$3]};}
    |   expression  '>' expression
        {$$ = {T:'bop',V:'>',C:[$1,$3]};}
    |   expression  '<' expression
        {$$ = {T:'bop',V:'<',C:[$1,$3]};}
    |   expression  '>=' expression
        {$$ = {T:'bop',V:'>=',C:[$1,$3]};}
    |   expression  '<=' expression
        {$$ = {T:'bop',V:'<=',C:[$1,$3]};}
    |   expression  '=' expression
        {$$ = {T:'bop',V:'=',C:[$1,$3]};}
    |   expression  '!=' expression
        {$$ = {T:'bop',V:'<>',C:[$1,$3]};}
    |   expression  'AND' expression
        {$$ = {T:'bop',V:'AND',C:[$1,$3]};}
    |   expression  'OR' expression
        {$$ = {T:'bop',V:'OR',C:[$1,$3]};}
    |   'NOT' expression
        {$$ = {T:'uop',V:'not',C:[$2]};}
    |   expression  'LIKE' expression
        {$$ = {T:'bop',V:'like',C:[$1,$3]};}
    ;

func_call
    :   ID '(' array_content ')'
        {$1.T=$1.T=='value'?'call_function':$1.T;$$={T:'function',V:$1,C:$3};}
    ;

array_content
    :   expression  ','  array_content
        {$$=$3; $$.unshift($1);}
    |   expression
        {$$=[$1];}
    |
        {$$=[];}
    ;

func_name
    :   ID
        {{$$=yytext;}}
    ;