import ast
import json
import sys

class CodebaseParser(ast.NodeVisitor):
    def __init__(self):
        self.entities = []
        self.relationships = []

    def visit_ClassDef(self, node):
        self.entities.append(('Class', node.name, node.lineno))
        for base in node.bases:
            if isinstance(base, ast.Name):
                self.relationships.append((node.name, 'INHERITS', base.id))
        self.generic_visit(node)

    def visit_FunctionDef(self, node):
        self.entities.append(('Function', node.name, node.lineno))
        for stmt in node.body:
            if isinstance(stmt, ast.Expr) and isinstance(stmt.value, ast.Call):
                if isinstance(stmt.value.func, ast.Name):
                    self.relationships.append((node.name, 'CALLS', stmt.value.func.id))
        self.generic_visit(node)

    def visit_Assign(self, node):
        if isinstance(node.targets[0], ast.Name):
            self.entities.append(('Variable', node.targets[0].id, node.lineno))
        self.generic_visit(node)

def parse_codebase(path):
    with open(path, "r") as source:
        tree = ast.parse(source.read())
    parser = CodebaseParser()
    parser.visit(tree)
    return parser.entities, parser.relationships

if __name__ == "__main__":
    path = sys.argv[1]
    entities, relationships = parse_codebase(path)
    result = {
        "Entities": [{"Label": e[0], "Name": e[1], "Line": e[2]} for e in entities],
        "Relationships": [{"From": r[0], "Type": r[1], "To": r[2]} for r in relationships]
    }
    print(json.dumps(result))
