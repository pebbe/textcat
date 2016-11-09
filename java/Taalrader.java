import nl.xs4all.pebbe.textcat.Textcat;

public class Taalrader {
    public static void main(String[] args) {
	String input = "Dit is een test in het Nederlands";
	if (args.length == 0) {
	    args = new String[1];
	    args[0] = "Dit is een test in het Nederlands";
	}
	Textcat textcat = new Textcat();
	textcat.setMinDocSize(10);
	for (int i = 0; i < args.length; i++) {
	    System.out.println(args[i]);
	    String result[] = textcat.classify(args[i]);
	    for (int j = 0; j < result.length; j++) {
		System.out.println(result[j]);
	    }
	}
    }
}
