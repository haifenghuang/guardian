
#include "ast.h"

CREATE_NODE(assignment)(struct expression_node** left, size_t token_position, struct expression_node** right);
CREATE_NODE(call)();
CREATE_NODE(return_stat)(struct expression_node** results);
CREATE_NODE(branch_stat)(const char* identifier);
CREATE_NODE(if_stat)(struct statement_node* init, struct expression_node* cond,
    struct block_stat_node* body, struct block_stat_node* else_branch);
CREATE_NODE(case_clause)(struct array* cases, struct block_stat_node* block);
CREATE_NODE(switch_stat)(struct node* switching_on, struct array* case_statements);
CREATE_NODE(block_stat)(struct node* list);
CREATE_NODE(for_stat)(struct statement_node* init, struct expression_node* cond,
    struct expression_node* post, struct block_stat_node block);

struct assignment_node {
    struct node base;
    struct expresion_node** left;
    struct expression_node** right;
};

// keyword <call>
// defer <call>
// run <call>
struct call_node {
    struct node base;
    struct call_expr_node* call;
};

struct return_stat_node {
    struct node base;
    struct expression_node** results;
};

struct branch_stat_node {
    struct node base;
    const char* identifier;
};

struct if_stat_node {
    struct node base;
    struct statement_node* init;
    struct expression_node* cond;
    struct block_stat_node* body;
    struct block_stat_node* else_branch;
};

struct case_clause_node {
    struct node base;
    struct array* cases; // expression nodes
    struct block_stat_node* block;
};

struct switch_stat_node {
    struct node base;
    struct expression_node* switching_on;
    struct array* case_statements;
};

struct block_stat_node {
    struct node base;
    struct array* statements;
};

struct for_stat_node {
    struct node base;
    struct statement_node* init;
    struct expression_node* cond;
    struct statement_node* post;
    struct block_stat_node* block;
};